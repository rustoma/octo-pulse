package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/sashabaranov/go-openai"
)

var logger *zerolog.Logger

type ChatGPTer interface {
	GenerateArticleDescription(question *models.Question) (string, error)
	AssignToCategory(categories []*models.Category, question *models.Question) (int, error)
	CheckIfPageContentIsValid(text string) (bool, error)
	CheckIfResponseContainRejected(response string) bool
	GenerateImage() (openai.ImageResponse, error)
}

type chatGPT struct {
	retriesLimit int
	Client       *openai.Client
	Usage        *openai.Usage
}

func NewChatGPT() ChatGPTer {
	client := openai.NewClient(os.Getenv("AI_KEY"))

	return &chatGPT{
		retriesLimit: 1,
		Client:       client,
		Usage:        &openai.Usage{PromptTokens: 0, CompletionTokens: 0, TotalTokens: 0},
	}
}

func (c *chatGPT) GenerateImage() (openai.ImageResponse, error) {

	resp, err := c.Client.CreateImage(context.Background(), openai.ImageRequest{
		Model:          openai.CreateImageModelDallE3,
		Prompt:         "Generate a person applying eye drops in the bathroom in a pleasant atmosphere. The person should have his eyes open and hit the eye drops.",
		Size:           "1024x1024",
		Quality:        "standard",
		N:              1,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
	})

	return resp, err
}

func (c *chatGPT) CorrectGrammar(text string) (string, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Zwróć bezpośrednio poprawiony tekst bez żadnego dodatkowego opisu. \n" +
				"Popraw gramatykę, błędy stylystyczne, składniowe oraz składnię HTML. Nie zmieniaj tekstu, ani struktury HTML, jedynie popraw błędy. \n" +
				"Tekst do poprawy: \n\n" +
				text,
		},
	}
	logger.Info().Msg("Starting to correct the text...")
	correctedText, err := c.ask(messages, openai.GPT3Dot5Turbo16K)
	if err != nil {
		logger.Error().Msgf("The text could not be corrected: %s \n\n Error: %v\n", text, err)
		return "", err
	}
	logger.Info().Msg("Text corrected successfully!")
	return correctedText, err
}

func (c *chatGPT) AssignToCategory(categories []*models.Category, question *models.Question) (int, error) {

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return 0, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Tytuł artykułu to: " + question.Question + "\n" +
				"Opis artykułu: " + question.Answer + "\n\n" +
				"Dostępne kategorie: " + string(categoriesJSON) + "\n\n" +
				"Przypasuj tytuł artykułu do jednej z podanych kategorii. Zwróć jedynie id kategorii. \n\n" +
				"Odpowiedź według zaleceń: \n\n" +
				"- zwróć jedynie id kategorii do której pasuje tytuł \n" +
				"- jeżeli tytuł nie pasuje do żadnej kategorii zwróć 0 \n" +
				"- id kategorii zwróć pomiędzy trzema myślnikami \n\n" +
				"Przykład poprawnej odpowiedzi: ---133---",
		},
	}

	resp, err := c.ask(messages)

	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`---(\d+)---`)
	match := re.FindStringSubmatch(resp)
	if match == nil {
		return 0, errors.New("No integer found in the respond from AssignToCategory")
	}

	categoryId, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return categoryId, nil
}

type Subtitle struct {
	Title     string   `json:"title"`
	Subtitles []string `json:"subtitles"`
}

type ArticleAgenda struct {
	MainTitle string     `json:"mainTitle"`
	Subtitles []Subtitle `json:"subtitles"`
}

func (c *chatGPT) CheckIfResponseContainRejected(response string) bool {
	// Regular expression pattern
	pattern := "---reject---"

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the string contains the substring
	if re.MatchString(response) {
		return true
	}

	return false
}

func (c *chatGPT) CheckIfPageContentIsValid(text string) (bool, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: "Przeanalizuj poniższy tekst i zwróć pomiędzy trzema myślnikami ---reject--- jeżeli: \n\n" +
				"1. Jeżeli tekst jest w języku angielskim. \n" +
				"2. Jeżeli 1 punkt nie pasuje to zwróć pomiędzy trzema myślnikami ---approved---",
		},
	}

	resp, err := c.ask(messages)
	if err != nil {
		return false, err
	}

	if c.CheckIfResponseContainRejected(resp) {
		return false, nil
	}

	return true, nil
}

func (c *chatGPT) RemoveMultipleSpaces(text string) string {
	trimmedText := strings.TrimSpace(text)

	// Define the regular expression pattern to match consecutive whitespaces
	pattern := `\s+`

	// Compile the regular expression pattern
	reg := regexp.MustCompile(pattern)

	// Replace multiple whitespaces with a single space
	return reg.ReplaceAllString(trimmedText, " ") + "\n\n"
}

func (c *chatGPT) GenerateArticleDescription(question *models.Question) (string, error) {

	var articleDescription bytes.Buffer

	var sourceText string

	for _, pageContent := range question.PageContents {
		text := c.RemoveMultipleSpaces(pageContent.PageContentProcessed)
		sourceText = fmt.Sprintf("%s \n\n %s", sourceText, text)
	}

	if len(sourceText) < 2000 {
		logger.Info().Msgf("There is no enough page content for the question. Question id: %d", question.Id)
		return "", nil
	}

	var messages []openai.ChatCompletionMessage

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: sourceText,
	})

	messages = append(messages,
		openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: "Głównym nagłówkiem będzie: " + question.Question + "\n Na podstawie tekstu który podałeś zwróć obiekt z nagłówkami i podrzędnymi nagłówkami, które mogą posłużyć do napisania takiego artykułu. \n" +
				"Nie uwzględniaj treści związanych z informacjami na temat firmy, polityki prywatności lub ciasteczek cookies. \n" +
				"Nie uwzględniaj tytułów takich jak 'o nas', 'informacje kontaktowe', 'o firmie' i wszystkich innych powiązanych z konkretną firmą. Nie dodawaj tytułów związanych z newsletter, biuletynem itp. \n" +
				"Nie dodawaj podtytułów 'podsumowanie', 'kontynuacja tematu' \n" +
				"Nie numeruj tytułów. \n" +
				"Zwróć poprawny json string na wzór: \n\n" +
				"{\"mainTitle:\"" + question.Question + ", \"subtitles\": [{\"title\": \"Subtitle1\", \"subtitles\": [\"Subtitle1\", \"Subtitle2\"]},{\"title\": \"Subtitle2\", \"subtitles\": [\"Subtitle1\", \"Subtitle2\"]}]}" + "\n\n" +
				"Nie dodawaj znaczników '\n'. Wszystko zwróć w jedej linii",
		})

	agenda, err := c.ask(messages)
	if err != nil {
		return "", err
	}

	var articleAgenda ArticleAgenda
	err = json.Unmarshal([]byte(agenda), &articleAgenda)
	if err != nil {
		return "", err
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: agenda,
	})

	logger.Info().Interface("agenda: ", articleAgenda).Send()

	var summary string

	logger.Info().Msg("Generating summary...")
	for index, pageContent := range question.PageContents {
		if len(c.RemoveMultipleSpaces(pageContent.PageContentProcessed)) < 1000 {
			continue
		}

		summaryPromp := []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleSystem,
			Content: pageContent.PageContentProcessed,
		},
			{
				Role: openai.ChatMessageRoleUser,
				Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z perfekcyjną znajomością języka polskiego. " +
					"Twoim celem jest stworzenie na podstawie tekstu, który podałeś podsumowania z najważniejszymi treścami pisanego jakby był to nowy artykuł, który będzie wykrozystany jako kontekst przy pisaniu artykułu do którego spis treści wygląda następująco: " + agenda,
			},
		}

		summaryResponse, err := c.ask(summaryPromp)
		if err != nil {
			return "", err
		}

		if index == 0 {
			summary = fmt.Sprintf("%s", summaryResponse)
		} else {
			summary = fmt.Sprintf("%s \n\n %s", summary, summaryResponse)
		}
	}
	logger.Info().Msg("Generating summary COMPLETED!")

	messages[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Streszczenie: " + summary,
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z perfekcyjną znajomością języka polskiego. " +
			"Twoim celem jest stworzyć 100% oryginalny, zoptymalizowany pod względem SEO artykuł, który czyta się jak napisany przez człowieka. " +
			"Styl odpowiedzi powinien być profesjonalny. Będzie to artykuł gdzie odbiorca będzie mógł zaczerpnąć informacji. \n\n" +
			"Na podstawie zadanego tytułu zwróć krótki wstęp do artykułu. \n\n" +
			"Podtytułami dla tego artykułu będą podtytuły jak w poniższej tablicy: \n\n" +
			fmt.Sprintf("%+v", articleAgenda.Subtitles) + "\n\n" +
			"Tytuł artykułu to: " + question.Question + "\n\n" +
			"Stosuj się do poniższych wymagań: \n\n" +
			"- Napisz tylko wstęp dla tego tytułu nie odpowiadaj na żadne podtytuły. \n" +
			"- Nie powtarzaj się. \n" +
			"- Nie pisz nic na temat SEO \n" +
			"- Długość wstępu powinina mieć minimum 1000 liter. \n" +
			"- Możesz zdefiniować kilka paragrafów, aby osiągnąć wymaganą długość wstępu. \n" +
			"- Długość wstępu jest wymagana! Powinna być bezwględnie przestrzegana! \n" +
			"- Tekst zwróć jako HTML. Tytuł artykułu powinien być w tagu <h1> \n" +
			"- Zwróć jedynie HTML z tekstem tak, aby dało się go dołączyć do już isntniejącego HTML. \n" +
			"- Odpowiedz jedynie HTML, tak abym całą odpowiedź mógł to skopiować i wkleić. \n" +
			"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <strong>, <h1> \n" +
			"- Nie dodawaj żadnych instrukcji od siebie.",
	})

	entryText, err := c.ask(messages)
	if err != nil {
		return "", err
	}

	correctedEntryText, err := c.CorrectGrammar(entryText)
	if err != nil {
		return "", err
	}

	articleDescription.WriteString(correctedEntryText)

	for _, subtitle := range articleAgenda.Subtitles {
		messagesLvl2 := []openai.ChatCompletionMessage{
			messages[0],
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: correctedEntryText,
			},
			{
				Role: openai.ChatMessageRoleUser,
				Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z perfekcyjną znajomością języka polskiego. " +
					"Twoim celem jest stworzyć 100% oryginalny, zoptymalizowany pod względem SEO artykuł, który czyta się jak napisany przez człowieka. " +
					"Styl odpowiedzi powinien być profesjonalny. Będzie to artykuł gdzie odbiorca będzie mógł zaczerpnąć informacji. \n\n" +
					"Rozwiń zadany podtytuł. \n\n" +
					"Podtytułami dla zadanego podtytułu będą podtytuły jak w poniższej tablicy: \n\n" +
					fmt.Sprintf("%+v", subtitle.Subtitles) + "\n\n" +
					"Zadany podtytuł to: " + subtitle.Title + "\n\n" +
					"Stosuj się do poniższych wymagań: \n\n" +
					"- Tekst zwróć jako HTML. Zadany podtytuł powinien być w tagu <h2> \n" +
					"- Zwróć jedynie HTML z tekstem tak, aby dało się go dołączyć do już isntniejącego HTML. \n" +
					"- Odpowiedz jedynie HTML, tak abym całą odpowiedź mógł to skopiować i wkleić. \n" +
					"- Odpowiedz jedynie za pomocą HTML. Nie pisz mi nic co mam z nim zrobić, ani że jest to odpowiedź." +
					"- W tekście nie odpowiadaj na żadne podtytuły. \n" +
					"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <strong>, <h2> \n" +
					"- Tekst powinien być powiązany kontekstem z głównym tytułem artykułu. \n" +
					"- Nie używaj w odpowiedzi tytułu nadrzędnego lub podtytułów dla zadanego podtytułu \n" +
					"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
					"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
					"- Nie powtarzaj się \n" +
					"- Nie pisz nic na temat SEO \n" +
					"- Możesz bazować na informacjach zawartych w streszczeniu. \n" +
					"- Jeżeli nie możesz udzielić lub kontynuować odpowiedzi zwróć pomiędzy trzema myślnikami ---reject--- \nn" +
					"Przykład poprawnej struktury odpowiedzi: \n\n" +
					"<h2>" + subtitle.Title + "</h2>" + "<p>...</p>",
			},
		}

		respLvl2, err := c.ask(messagesLvl2)
		if err != nil {
			return "", err
		}

		if c.CheckIfResponseContainRejected(respLvl2) {
			logger.Info().Msgf("Text cannot be proccess %s", subtitle.Title)
			continue
		}

		correctedRespLvl2, err := c.CorrectGrammar(respLvl2)
		if err != nil {
			return "", err
		}

		articleDescription.WriteString(correctedRespLvl2)

		logger.Info().Interface("subtitle LVL2: ", subtitle.Title).Send()

		var allMessagesLvl3 []openai.ChatCompletionMessage
		for index, subtitle3lvl := range subtitle.Subtitles {
			if index == 1 {
				messagesLvl3 := []openai.ChatCompletionMessage{
					messagesLvl2[1],
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: correctedRespLvl2,
					},
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z perfekcyjną znajomością języka polskiego. " +
							"Twoim celem jest stworzyć 100% oryginalny, zoptymalizowany pod względem SEO artykuł, który czyta się jak napisany przez człowieka. " +
							"Styl odpowiedzi powinien być profesjonalny. Będzie to artykuł gdzie odbiorca będzie mógł zaczerpnąć informacji. \n\n" +
							"Rozwiń zadany podtytuł. \n\n" +
							"Zadany tytuł jest to podtytuł tytułu nadrzędnego jak poniżej: \n\n" +
							subtitle.Title + "\n\n" +
							"Zadany podtytuł to: " + subtitle3lvl + "\n\n" +
							"Stosuj się do poniższych wymagań: \n\n" +
							"- Tekst zwróć jako HTML. Zadany podtytuł powinien być w tagu <h3> \n" +
							"- Zwróć jedynie HTML z tekstem tak, aby dało się go dołączyć do już isntniejącego HTML. \n" +
							"- Odpowiedz jedynie HTML, tak abym całą odpowiedź mógł to skopiować i wkleić. \n" +
							"- Odpowiedz jedynie za pomocą HTML. Nie pisz mi nic co mam z nim zrobić, ani że jest to odpowiedź." +
							"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <strong>, <h3> \n" +
							"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
							"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
							"- Nie powtarzaj się \n" +
							"- Nie pisz nic na temat SEO \n" +
							"- Możesz bazować na informacjach zawartych w streszczeniu. \n" +
							"- Jeżeli nie możesz udzielić lub kontynuować odpowiedzi zwróć pomiędzy trzema myślnikami ---reject--- \nn" +
							"Przykład poprawnej struktury odpowiedzi: \n\n" +
							"<h3>" + subtitle3lvl + "</h3>" + "<p>...</p>",
					},
				}
				allMessagesLvl3 = messagesLvl3
			} else {
				messagesLvl3 := []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z perfekcyjną znajomością języka polskiego. " +
							"Twoim celem jest stworzyć 100% oryginalny, zoptymalizowany pod względem SEO artykuł, który czyta się jak napisany przez człowieka. " +
							"Styl odpowiedzi powinien być profesjonalny. Będzie to artykuł gdzie odbiorca będzie mógł zaczerpnąć informacji. \n\n" +
							"Rozwiń zadany podtytuł. \n\n" +
							"Zadany tytuł jest to podtytuł tytułu nadrzędnego jak poniżej: \n\n" +
							subtitle.Title + "\n\n" +
							"Zadany podtytuł to: " + subtitle3lvl + "\n\n" +
							"Stosuj się do poniższych wymagań: \n\n" +
							"- Tekst zwróć jako HTML. Zadany podtytuł powinien być w tagu <h3> \n" +
							"- Zwróć jedynie HTML z tekstem tak, aby dało się go dołączyć do już isntniejącego HTML. \n" +
							"- Odpowiedz jedynie HTML, tak abym całą odpowiedź mógł to skopiować i wkleić. \n" +
							"- Odpowiedz jedynie za pomocą HTML. Nie pisz mi nic co mam z nim zrobić, ani że jest to odpowiedź." +
							"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <strong>, <h3> \n" +
							"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
							"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
							"- Nie powtarzaj się \n" +
							"- Nie pisz nic na temat SEO \n" +
							"- Możesz bazować na informacjach zawartych w streszczeniu. \n" +
							"- Jeżeli nie możesz udzielić lub kontynuować odpowiedzi zwróć pomiędzy trzema myślnikami ---reject--- \nn" +
							"Przykład poprawnej struktury odpowiedzi: \n\n" +
							"<h3>" + subtitle3lvl + "</h3>" + "<p>...</p>",
					},
				}
				allMessagesLvl3 = append(allMessagesLvl3, messagesLvl3...)
			}

			respLvl3, err := c.ask(allMessagesLvl3)
			if err != nil {
				return "", err
			}

			if c.CheckIfResponseContainRejected(respLvl3) {
				logger.Info().Msgf("Text cannot be proccess %s", subtitle3lvl)
				continue
			}

			correctedRespLvl3, err := c.CorrectGrammar(respLvl3)
			if err != nil {
				return "", err
			}

			articleDescription.WriteString(correctedRespLvl3)

			logger.Info().Interface("subtitle 3lvl: ", subtitle3lvl).Send()
		}
	}

	//promptTokenCost := float64(c.Usage.PromptTokens) / 1000 * 0.01
	//completionTokensCost := float64(c.Usage.CompletionTokens) / 1000 * 0.03
	//totalTokenCost := promptTokenCost + completionTokensCost
	//
	//promptTokenCostInDolars := fmt.Sprintf("%f $", promptTokenCost)
	//completionTokensCostInDolars := fmt.Sprintf("%f $", completionTokensCost)
	//totalTokenCostInDolars := fmt.Sprintf("%f $", totalTokenCost)
	//
	//articleDescription.WriteString("<p>Usage:</p>" +
	//	"<p>PromptTokenCost: " + fmt.Sprintf("%d", c.Usage.PromptTokens) + " - " + promptTokenCostInDolars + "</p>" +
	//	"<p>CompletionTokenCost: " + fmt.Sprintf("%d", c.Usage.CompletionTokens) + " - " + completionTokensCostInDolars + "</p>" +
	//	"<p>TotalTokenCost: " + fmt.Sprintf("%d", c.Usage.TotalTokens) + " - " + totalTokenCostInDolars + "</p>")

	c.Usage = &openai.Usage{PromptTokens: 0, CompletionTokens: 0, TotalTokens: 0}

	re := regexp.MustCompile("<h1[^>]*>(.*?)</h1>")
	articleDescriptionWithoutH1 := re.ReplaceAllString(articleDescription.String(), "")

	return articleDescriptionWithoutH1, nil
}

func (c *chatGPT) newChatCompletion(messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {

	logger.Info().Msg(model)

	resp, err := c.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: 0,
		},
	)

	return resp, err
}

func (c *chatGPT) retry(messages []openai.ChatCompletionMessage, retriesLimit int, model string) (openai.ChatCompletionResponse, error) {
	var (
		retries     = 0
		isAskFailed = true
	)

	for {
		logger.Info().Msg("Re-attempting...")
		retries += 1
		logger.Info().Msgf("Retry number: %d", retries)

		resp, err := c.newChatCompletion(messages, model)

		if err != nil {
			logger.Err(err).Send()
		} else {
			isAskFailed = false
		}

		if retries > retriesLimit || !isAskFailed {
			logger.Info().Msg("Repeated query successful")
			return resp, err
		}
	}

}

func (c *chatGPT) ask(messages []openai.ChatCompletionMessage, model ...string) (string, error) {
	chatCompletionModel := openai.GPT4TurboPreview
	if len(model) > 0 {
		chatCompletionModel = model[0]
	}

	resp, err := c.newChatCompletion(messages, chatCompletionModel)

	if err != nil {
		logger.Err(err).Send()
		resp, err = c.retry(messages, c.retriesLimit, chatCompletionModel)
	}

	if err != nil {
		return "", err
	}

	c.Usage = &openai.Usage{
		PromptTokens:     c.Usage.PromptTokens + resp.Usage.PromptTokens,
		CompletionTokens: c.Usage.CompletionTokens + resp.Usage.CompletionTokens,
		TotalTokens:      c.Usage.TotalTokens + resp.Usage.TotalTokens,
	}

	logger.Info().Interface("Usage: ", resp.Usage).Send()
	return resp.Choices[0].Message.Content, nil
}

func init() {
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
