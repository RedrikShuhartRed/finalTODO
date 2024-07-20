package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	environment "github.com/RedrikShuhartRed/finalTODO/Environment"
	"github.com/golang-jwt/jwt"
)

func HashPassword(password string) string {
	passwordBytes := []byte(environment.LoadEnvPassword())
	passwordSaltBytes := []byte(environment.LoadEnvPasswordSalt())
	passwordBytes = append(passwordBytes, passwordSaltBytes...)
	hashedPasswordBytes := sha256.Sum256(passwordBytes)
	// claims := jwt.MapClaims{
	// 	"hashPass": hex.EncodeToString(hashedPasswordBytes[:]),
	// }
	hashedPasswordBytesHex := hex.EncodeToString(hashedPasswordBytes[:])
	return hashedPasswordBytesHex
}

func GenerateJWT(hashedPasswordBytesHex string) (string, error) {
	claims := jwt.MapClaims{
		"hashPass": hex.EncodeToString([]byte(hashedPasswordBytesHex[:])),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSaltBytes := []byte(environment.LoadEnvTokenSalt())

	signedToken, err := jwtToken.SignedString(tokenSaltBytes)
	if err != nil {
		log.Printf("error , %v", err)
		return "", err
	}
	return signedToken, nil
}
