package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSONError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequestHandler(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusBadRequest, "bad_request", message)
}

func ValidationErrorHandler(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusBadRequest, "validation_error", message)
}

// Функция дла обработки ошибки, связанной с некорректной работой сервера.
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}

// Функция для обработки ошибки
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotFound, "not_found", "not found")
}
