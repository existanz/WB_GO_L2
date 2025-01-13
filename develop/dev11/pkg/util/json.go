package util

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func WriteResult(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(map[string]string{"result": message})
}

func WriteJSON(w http.ResponseWriter, v any) {
	message, err := json.Marshal(v)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResult(w, string(message))
}

func BindJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
