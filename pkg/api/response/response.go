package response

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, err error, statusCode int) {
	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = http.StatusText(statusCode)
	}

	JSON(w, &errorMessage{Message: message}, statusCode)
}

func NotFound(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusNotFound)
}

func InternalError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusInternalServerError)
}

func JSON(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
