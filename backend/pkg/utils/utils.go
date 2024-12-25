package utils

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	Response(w, r, code, map[string]string{"error": err.Error()})
}

func Response(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
