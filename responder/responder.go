package responder

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Error responds with a HTTP status code of 500 and writes the error message to the body
func Error(w http.ResponseWriter, err error) error {
	return ErrorWithStatus(http.StatusInternalServerError, w, err)
}

// ErrorWithStatus responds with an error
func ErrorWithStatus(status int, w http.ResponseWriter, err error) error {
	return JSON(status, w, map[string]interface{}{
		"ErrorTitle":   http.StatusText(status),
		"ErrorMessage": err.Error(),
	})
}

// JSON writes the JSON representation of of the object into the HTTP response
func JSON(status int, w http.ResponseWriter, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}

	w.WriteHeader(status)

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}
	if w.Header().Get("Content-Length") == "" {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	}

	w.Write(data)

	return nil
}
