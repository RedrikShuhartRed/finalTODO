package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonErr struct {
	err error
}

func (err *JsonErr) JsonError(w http.ResponseWriter, message string, code int) {
	resp := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errEncode := json.NewEncoder(w).Encode(resp)
	if errEncode != nil {
		log.Printf("error encode error resrp, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
