package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hashString string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hashString)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStructs := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStructs, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("there was was an error while trying to decode the token: %w", err)
	}

	subject, issueError := token.Claims.GetSubject()

	if issueError != nil {
		return uuid.Nil, fmt.Errorf("there was a problem while trying to get the subject: %w", issueError)
	}

	subjectParsed, err := uuid.Parse(subject)
	return subjectParsed, err
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	authSlice := strings.Fields(authHeader)

	if len(authSlice) == 0 {
		return "", errors.New("there is no authorization header in the request")
	}

	if authSlice[0] != "Bearer" || len(authSlice[1]) == 0 {
		return "", fmt.Errorf("the authorization Header is not valid: %s", authHeader)
	}

	return authSlice[1], nil
}
