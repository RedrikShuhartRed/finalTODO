package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/RedrikShuhartRed/finalTODO/api/middleware"
)

var (
	errWrongPassword = errors.New("error: wrong password")
	errTokenPassword = errors.New("error: token password hash doesn't match")
)

func (h *Handler) AuthorizationGetToken(w http.ResponseWriter, r *http.Request) {
	password := map[string]string{}

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if password["password"] != h.cfg.Password {
		log.Printf("error: wrong password %v", errWrongPassword)
		Jerr.JsonError(w, errWrongPassword.Error(), http.StatusUnauthorized)

		return
	}
	hashedPasswordBytesHex := middleware.HashPassword(password["password"], h.cfg)
	signedToken, err := middleware.GenerateJWT(hashedPasswordBytesHex, h.cfg)
	if err != nil {
		log.Printf("error , %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(&map[string]string{"token": signedToken})
	_, err = w.Write(res)
	if err != nil {
		log.Printf("error during writing data to response writer %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(h.cfg.Password) > 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			hashPass, err := middleware.GetHashFromCockie(r, h.cfg)
			if err != nil {
				log.Printf("error , %v", err)
				Jerr.JsonError(w, err.Error(), http.StatusUnauthorized)
			}

			hashedPasswordBytesHex := middleware.HashPassword(h.cfg.Password, h.cfg)

			if hashPass != hashedPasswordBytesHex {
				log.Printf("error, %v", errTokenPassword)
				Jerr.JsonError(w, errTokenPassword.Error(), http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
