package handler

import (
	"fmt"
	//"os"
	"net/http"
)

// Ручка для домашней страницы.
type HomeHandler struct{}

// Функция для обработки запроса.
func (h *HomeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("This page do nothing."))
}

// NOTE: Я не знаю, где это еще написать.
func serveCSV(w http.ResponseWriter, r *http.Request, filename string) {
    w.Header().Set("Content-Type", "text/csv; charset=utf-8")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    
    http.ServeFile(w, r, filename)
    
    // Удаляем временный файл после отправки, чтобы не засорять диск
	// NOTE: я, пожалуй, это пока оставлю. В будущем надо будет убрать и вообще с нормальной бд это запускать.
    //os.Remove(filename) 
}
