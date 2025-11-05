package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var emailRequest requestEmail

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&emailRequest); err != nil {
		jsonResponseError(w, 403, fmt.Sprintf("there was a problem while decoding the request: %v", err))
		return
	}

	newUser, err := cfg.db.CreateUser(context.Background(), emailRequest.Email)

	if err != nil {
		jsonResponseError(w, 403, fmt.Sprintf("there was a problem while decoding the request: %v", err))
	}

	emailRes := User{
		Id:         newUser.ID,
		Created_at: newUser.CreatedAt,
		Updated_at: newUser.UpdatedAt,
		Email:      newUser.Email,
	}

	jsonResponse(w, 201, emailRes)
}
