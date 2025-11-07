package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Raghart/go-servers/internal/auth"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		jsonResponseError(w, 401, "There is no refresh token in the header")
		return
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(context.Background(), refreshToken)

	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("there is no refresh token that match the requested token: %v", err))
		return
	}

	acessToken, err := auth.MakeJWT(dbRefreshToken.ID, cfg.secretString, time.Duration(3600)*time.Second)

	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("there was a problem while trying to create the acess token: %v", err))
	}

	jsonResponse(w, 200, TokenResponse{
		Token: acessToken,
	})
}
