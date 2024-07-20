package middleware

import (
	"encoding/json"
	"net/http"
)

type JsonErr struct {
	err error
}

func (err *JsonErr) JsonError(w http.ResponseWriter, message string, code int) {
	resp := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
