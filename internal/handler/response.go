package handler

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

func writeJSONResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(SuccessResponse{Data: payload})
}
