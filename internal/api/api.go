package api

import (
	"encoding/json"
	"errors"
	e "github.com/rustoma/octo-pulse/internal/errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Err    string
	Status int
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (e Error) Error() string {
	return e.Err
}

func HandleErrorStatus(err interface{}) int {
	switch err.(type) {
	case e.Unauthorized:
		return http.StatusUnauthorized
	case e.BadRequest, e.NotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 //one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()
	return WriteJSON(w, statusCode, payload)
}

func MakeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(Error); ok {
				_ = WriteJSON(w, e.Status, e)
				return
			}
			_ = WriteJSON(w, http.StatusInternalServerError, Error{Err: "Internal sever", Status: http.StatusInternalServerError})
		}
	}
}
