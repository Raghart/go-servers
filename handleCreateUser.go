package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Raghart/go-servers/internal/auth"
	"github.com/Raghart/go-servers/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var emailRequest requestUser

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&emailRequest); err != nil {
		jsonResponseError(w, 403, fmt.Sprintf("there was a problem while decoding the request: %v", err))
		return
	}

	hashedPassword, err := auth.HashPassword(emailRequest.Password)

	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("the password isn't in a valid format: %v", err))
	}

	userParams := database.CreateUserParams{
		Email:          emailRequest.Email,
		HashedPassword: hashedPassword,
	}

	newUser, err := cfg.db.CreateUser(context.Background(), userParams)

	if err != nil {
		jsonResponseError(w, 403, fmt.Sprintf("there was a problem while decoding the request: %v", err))
	}

	emailRes := User{
		Id:            newUser.ID,
		Created_at:    newUser.CreatedAt,
		Updated_at:    newUser.UpdatedAt,
		Email:         newUser.Email,
		Is_chirpy_red: newUser.IsChirpyRed,
	}

	jsonResponse(w, 201, emailRes)
}
