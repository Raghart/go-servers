package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Raghart/go-servers/internal/auth"
	"github.com/Raghart/go-servers/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var userRequest requestUser

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&userRequest); err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("there was a problem while trying to decode the request: %v", err))
		return
	}

	databaseUser, err := cfg.db.GetUserByEmail(context.Background(), userRequest.Email)

	if err != nil {
		jsonResponseError(w, 401, "Incorrect email or password")
		return
	}

	isMatch, err := auth.CheckPasswordHash(userRequest.Password, databaseUser.HashedPassword)

	if err != nil || !isMatch {
		jsonResponseError(w, 401, "Incorrect email or password")
		return
	}

	acessToken, err := auth.MakeJWT(databaseUser.ID, cfg.secretString, time.Duration(3600)*time.Second)

	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("there was a problem while trying to make the acess token: %v", err))
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		jsonResponseError(w, 500, fmt.Sprintf("Error while creating the refresh token: %v", err))
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    databaseUser.ID,
		ExpiresAt: time.Now().Add(time.Duration(5184000) * time.Second),
	}

	newRefreshToken, err := cfg.db.CreateRefreshToken(context.Background(), refreshTokenParams)

	if err != nil {
		jsonResponseError(w, 500, fmt.Sprintf("error while trying to create the refresh Token: %v", err))
		return
	}

	jsonResponse(w, 200, UserLogin{
		Id:            databaseUser.ID,
		Created_at:    databaseUser.CreatedAt,
		Updated_at:    databaseUser.UpdatedAt,
		Email:         databaseUser.Email,
		Token:         acessToken,
		Refresh_token: newRefreshToken.Token,
		Is_chirpy_red: databaseUser.IsChirpyRed,
	})
}
