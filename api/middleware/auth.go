package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/golang-jwt/jwt"
)

func HashPassword(password string, cfg *config.Config) string {
	passwordBytes := []byte(cfg.Password)
	passwordSaltBytes := []byte(cfg.PasswordSalt)
	passwordBytes = append(passwordBytes, passwordSaltBytes...)
	hashedPasswordBytes := sha256.Sum256(passwordBytes)
	hashedPasswordBytesHex := hex.EncodeToString(hashedPasswordBytes[:])
	return hashedPasswordBytesHex
}

func GenerateJWT(hashedPasswordBytesHex string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"hashPass": hashedPasswordBytesHex,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSaltBytes := []byte(cfg.TokenSalt)

	signedToken, err := jwtToken.SignedString(tokenSaltBytes)
	if err != nil {
		log.Printf("error , %v", err)
		return "", err
	}
	return signedToken, nil
}

func GetHashFromCockie(r *http.Request, cfg *config.Config) (string, error) {
	var token string
	cookie, err := r.Cookie("token")
	if err == nil {
		token = cookie.Value
	}

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.TokenSalt), nil
	})

	if err != nil {
		log.Printf("error , %v", err)
		return "", err
	}
	if !jwtToken.Valid {
		log.Printf("error jwt token isn't valid")
		return "", err
	}

	res, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("failed to typecast to jwt.MapCalims, %v", err)
		return "", err
	}

	hashPassRaw := res["hashPass"]
	hashPass, ok := hashPassRaw.(string)
	if !ok {
		log.Printf("failed to typecase password hash to string, %v", err)
		return "", err
	}
	return hashPass, nil
}
