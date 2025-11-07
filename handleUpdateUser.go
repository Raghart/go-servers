package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Raghart/go-servers/internal/auth"
	"github.com/Raghart/go-servers/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, req *http.Request) {
	acessToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		jsonResponseError(w, 401, "the acess token is missing or is malformatted")
		return
	}

	userID, err := auth.ValidateJWT(acessToken, cfg.secretString)

	if err != nil {
		jsonResponse(w, 401, fmt.Sprintf("validation error for the acess token: %v", err))
		return
	}

	defer req.Body.Close()
	var requestInfo requestUser

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&requestInfo); err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("there was a problem while decoding the request body: %v", err))
		return
	}

	hashedPassword, err := auth.HashPassword(requestInfo.Password)

	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("error trying to hash the password: %v", err))
		return
	}

	userParams := database.UpdateUserByEmailParams{
		Email:          requestInfo.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	userUpdated, err := cfg.db.UpdateUserByEmail(context.Background(), userParams)

	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("there was a problem while trying to update the user information: %v", err))
		return
	}

	jsonResponse(w, 200, User{
		Id:            userID,
		Created_at:    userUpdated.CreatedAt,
		Updated_at:    userUpdated.UpdatedAt,
		Email:         userUpdated.Email,
		Is_chirpy_red: userUpdated.IsChirpyRed,
	})
}
