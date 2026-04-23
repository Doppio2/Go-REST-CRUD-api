package handler

import (
	"fmt"
	"net/http"
)

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("This page do nothing."))
}

// serveCSV streams a generated CSV file back to the client.
func serveCSV(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	http.ServeFile(w, r, filename)
}
