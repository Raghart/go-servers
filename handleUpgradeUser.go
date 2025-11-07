package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Raghart/go-servers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUpgradeUser(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("api key not valid: %v", err))
		return
	}

	if apiKey != cfg.polkaKey {
		jsonResponseError(w, 401, "the provided API Key doesn't match with the polka key")
		return
	}

	defer req.Body.Close()
	var webhookData webHookRequest

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&webhookData); err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("there was a problem while trying to decode the body: %v", err))
		return
	}

	if webhookData.Event != "user.upgraded" {
		jsonResponse(w, 204, []byte("webhook request in not the type of user.upgraded!"))
	}

	userId, err := uuid.Parse(webhookData.Data.User_id)
	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("the given id is not a valid uuid format: %v", err))
	}

	_, err = cfg.db.UpgradeUserByID(context.Background(), userId)
	if err != nil {
		jsonResponseError(w, http.StatusNotFound, fmt.Sprintf("user not found in the database: %v", err))
	}

	jsonResponse(w, 204, []byte(""))
}
