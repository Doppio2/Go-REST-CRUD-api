package handler

import (
	"net/http"
)

// Ручка для домашней страницы.
type HomeHandler struct{}

// Функция для обработки запроса.
func (h *HomeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("This page do nothing."))
}
