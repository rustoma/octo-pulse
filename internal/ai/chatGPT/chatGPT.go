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

	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/sashabaranov/go-openai"
)

var logger *zerolog.Logger

type ChatGPTer interface {
	GenerateArticleDescription(question *models.Question) (string, error)
	AssignToCategory(categories []*models.Category, question *models.Question) (int, error)
}

type chatGPT struct {
	retriesLimit int
	Client       *openai.Client
}

func NewChatGPT() ChatGPTer {
	client := openai.NewClient(os.Getenv("AI_KEY"))

	return &chatGPT{
		retriesLimit: 2,
		Client:       client,
	}
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
				"Opis artykułu: " + question.Answear + "\n\n" +

				"Dostępne kategorie: " + string(categoriesJSON) + "\n\n" +

				"Przypasuj tytuł artykułu do jednej z podanych kategorii. Zwróć jedynie id kategorii. \n\n" +

				"Odpowiedź według zaleceń: \n\n" +

				"- zwróć jedynie id kategorii do której pasuje tytuł \n" +
				"- jeżeli tytuł nie pasuje do żadnej kategorii zwróć 0 \n" +
				"- id kategorii zwróć pomiędzy trzema myślnikami \n\n" +

				"Przykład poprawnej odpowiedzi: ---133---",
		},
	}

	logger.Info().Interface("message: ", messages).Send()

	resp, err := c.ask(messages)

	logger.Info().Interface("respond: ", resp).Send()
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

func (c *chatGPT) GenerateArticleDescription(question *models.Question) (string, error) {

	var articleDescription bytes.Buffer

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Napisz streszczenie, które może posłużyć jako informację do napisania artykułu po polsku. \n" +
				"Wypisz w punktach jakie tematy zostaną poruszone na podstawie poniższych danych: \n\n" +

				"Temat główny: " + question.Question + "\n\n" +

				question.Answear,
		},
	}

	resp, err := c.ask(messages)
	if err != nil {
		return "", err
	}

	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: resp,
		},
		openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: "Na podstawie streszczenia które podałeś zwróć obiekt z nagłówkami i podrzędnymi nagłówkami, które mogą posłużyć do napisania takiego artykułu.\n" +
				"Zwróć poprawny json string na wzór: \n\n" +

				`{"mainTitle": {{question.Question}}, "subtitles": [{"title": "Subtitle1", "subtitles": ["Subtitle1", "Subtitle2"]},{"title": "Subtitle2", "subtitles": ["Subtitle1", "Subtitle2"]}]}` + "\n\n" +
				"Nie dodawaj znaczników '\n'. Wszystko zwróć w jedej linii",
		})

	resp, err = c.ask(messages)
	if err != nil {
		return "", err
	}

	logger.Info().Interface("resp: ", resp).Send()

	var articleAgenda ArticleAgenda
	err = json.Unmarshal([]byte(resp), &articleAgenda)
	if err != nil {
		return "", err
	}

	logger.Info().Interface("articleAgenda: ", articleAgenda).Send()

	messages = []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z zaawansowaną wiedzą na temat SEO i perfekcyjną znajomością języka polskiego. " +
				"Twoim celem jest stworzyć 100% oryginalny, zoptymalizowany pod względem SEO artykuł, który czyta się jak napisany przez człowieka. " +
				"Styl odpowiedzi powinien być profesjonalny. Będzie to artykuł gdzie odbiorca będzie mógł zaczerpnąć informacji. \n\n" +
				"Na podstawie zadanego tytułu zwróć krótki wstęp do artykułu. \n\n" +
				"Podtytułami dla tego artykułu będą podtytuły jak w poniższej tablicy: \n\n" +
				fmt.Sprintf("%+v", articleAgenda.Subtitles) + "\n\n" +

				"Tytuł artykułu to: " + articleAgenda.MainTitle + "\n\n" +

				"Stosuj się do poniższych wymagań: \n\n" +

				"- Napisz tylko wstęp dla tego tytułu nie odpowiadaj na żadne podtytuły. \n" +
				"- Nie powtarzaj się. \n" +
				"- Długość wstępu powinina mieć minimum 1500 liter. \n" +
				"- Możesz zdefiniować kilka paragrafów, aby osiągnąć wymaganą długość wstępu. \n" +
				"- Długość wstępu jest wymagana! Powinna być bezwględnie przestrzegana! \n" +
				"- Wstęp powinnien zaweirać co najmniej 5 paragrafów. Paragraf to tag <p> \n" +
				"- Tekst zwróć jako HTML. Tytuł artykułu powinien być w tagu <h1> \n" +
				"- Zwróć jedynie HTML z tekstem tak, aby dało się go dołączyć do już isntniejącego HTML. \n" +
				"- Odpowiedz jedynie HTML, tak abym całą odpowiedź mógł to skopiować i wkleić. \n" +
				"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <br>, <strong>, <h1> \n" +
				"- Nie dodawaj żadnych instrukcji od siebie.",
		},
	}

	resp, err = c.ask(messages)
	if err != nil {
		return "", err
	}
	articleDescription.WriteString(resp)

	logger.Info().Interface("intro for article: ", resp).Send()

	for _, subtitle := range articleAgenda.Subtitles {
		messagesLvl2 := []openai.ChatCompletionMessage{
			messages[0],
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: resp,
			},
			{
				Role: openai.ChatMessageRoleUser,
				Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z zaawansowaną wiedzą na temat SEO i perfekcyjną znajomością języka polskiego. " +
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
					"- Tekst powinnien zawierać co najmniej 5 paragrafów. Paragraf to tag <p> \n" +
					"- Tekst powinnien być długi. \n" +
					"- W tekście nie odpowiadaj na żadne podtytuły. \n" +
					"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <br>, <strong>, <h2> \n" +
					"- Tekst powinien być powiązany kontekstem z głównym tytułem artykułu. \n" +
					"- Nie używaj w odpowiedzi tytułu nadrzędnego lub podtytułów dla zadanego podtytułu \n" +
					"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
					"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
					"- Nie powtarzaj się" +

					"Przykład poprawnej struktury odpowiedzi: \n\n" +
					"<h2>" + subtitle.Title + "</h2>" + "<p>...</p>",
			},
		}

		respLvl2, err := c.ask(messagesLvl2)
		if err != nil {
			return "", err
		}
		articleDescription.WriteString(respLvl2)

		logger.Info().Interface("subtitle: ", respLvl2).Send()

		var allMessagesLvl3 []openai.ChatCompletionMessage
		for index, subtitle3lvl := range subtitle.Subtitles {
			if index == 1 {
				messagesLvl3 := []openai.ChatCompletionMessage{
					messagesLvl2[1],
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: respLvl2,
					},
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z zaawansowaną wiedzą na temat SEO i perfekcyjną znajomością języka polskiego. " +
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
							"- Tekst powinnien być długi. \n" +
							"- Tekst powinnien zawierać co najmniej 5 paragrafów. Paragraf to tag <p> \n" +
							"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <br>, <strong>, <h3> \n" +
							"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
							"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
							"- Nie powtarzaj się" +

							"Przykład poprawnej struktury odpowiedzi: \n\n" +
							"<h3>" + subtitle3lvl + "</h3>" + "<p>...</p>",
					},
				}
				allMessagesLvl3 = messagesLvl3
			} else {
				messagesLvl3 := []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Wyobraź sobie, że jesteś doświadczonym copywriterem z zaawansowaną wiedzą na temat SEO i perfekcyjną znajomością języka polskiego. " +
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
							"- Tekst powinnien być długi. \n" +
							"- Tekst powinnien zawierać co najmniej 5 paragrafów. Paragraf to tag <p> \n" +
							"- Dozwolone tagi HTML to : <p>, <ul>, <li>, <ol>, <br>, <strong>, <h3> \n" +
							"- Tekst powinien być powiązany kontekstem z poprzednimi odpowiedziami. \n" +
							"- Nie używaj w odpowiedzi tytułu nadrzędnego \n" +
							"- Nie powtarzaj się" +

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
			articleDescription.WriteString(respLvl3)

			logger.Info().Interface("subtitle 3lvl: ", respLvl3).Send()
		}

	}

	return articleDescription.String(), nil
}

func (ai *chatGPT) newChatCompletion(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {

	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K,
			Messages:    messages,
			Temperature: 0,
		},
	)

	return resp, err
}

func (ai *chatGPT) retry(messages []openai.ChatCompletionMessage, retriesLimit int) (openai.ChatCompletionResponse, error) {
	var (
		retries     = 0
		isAskFailed = true
	)

	for {
		logger.Info().Msg("Re-attempting...")
		retries += 1
		logger.Info().Msgf("Retry number: %d", retries)

		resp, err := ai.newChatCompletion(messages)

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

func (ai *chatGPT) ask(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := ai.newChatCompletion(messages)

	if err != nil {
		logger.Err(err).Send()
		resp, err = ai.retry(messages, ai.retriesLimit)
	}

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func init() {
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
