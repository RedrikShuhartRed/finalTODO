package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	environment "github.com/RedrikShuhartRed/finalTODO/Environment"
	"github.com/RedrikShuhartRed/finalTODO/api/middleware"
	"github.com/golang-jwt/jwt"
)

func AuthorizationGetToken(w http.ResponseWriter, r *http.Request) {
	password := map[string]string{}

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		jsonError(w, "error Decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if password["password"] != environment.LoadEnvPassword() {
		log.Printf("error wrong password")
		jsonError(w, "error wrong password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// passwordBytes := []byte(environment.LoadEnvPassword())
	// passwordSaltBytes := []byte(environment.LoadEnvPasswordSalt())
	// passwordBytes = append(passwordBytes, passwordSaltBytes...)
	// hashedPasswordBytes := sha256.Sum256(passwordBytes)
	hashedPasswordBytesHex := middleware.HashPassword(password["password"])
	signedToken, err := middleware.GenerateJWT(hashedPasswordBytesHex)
	if err != nil {
		log.Printf("error , %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// claims := jwt.MapClaims{
	// 	"hashPass": hex.EncodeToString([]byte(hashedPasswordBytesHex[:])),
	// }
	// jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenSaltBytes := []byte(environment.LoadEnvTokenSalt())

	// signedToken, err := jwtToken.SignedString(tokenSaltBytes)
	// if err != nil {
	// 	log.Printf("error , %v", err)
	// 	jsonError(w, err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(&map[string]string{"token": signedToken})
	_, err = w.Write(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during writing data to response writer %s", err.Error())
		return
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(environment.LoadEnvPassword()) > 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}

			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(environment.LoadEnvTokenSalt()), nil
			})

			if err != nil {
				log.Printf("error , %v", err)
				jsonError(w, err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !jwtToken.Valid {
				log.Printf("error jwt token isn't valid")
				jsonError(w, "jwt token isn't valid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				log.Printf("failed to typecast to jwt.MapCalims, %v", err)
				jsonError(w, "failed to typecast to jwt.MapCalims")
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			hashPassRaw := res["hashPass"]
			hashPass, ok := hashPassRaw.(string)
			if !ok {
				log.Printf("failed to typecase password hash to string, %v", err)
				jsonError(w, "failed to typecase password hash to string")
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			passwordBytes := []byte(environment.LoadEnvPassword())
			passwordSaltBytes := []byte(environment.LoadEnvPasswordSalt())
			passwordBytes = append(passwordBytes, passwordSaltBytes...)
			hashedPasswordBytes := sha256.Sum256(passwordBytes)
			if hashPass != hex.EncodeToString(hashedPasswordBytes[:]) {
				log.Printf("token password hash doesn't match, %v", err)
				jsonError(w, "token password hash doesn't match")
				w.WriteHeader(http.StatusUnauthorized)

				return
			}
		}
		next(w, r)
	})
}
