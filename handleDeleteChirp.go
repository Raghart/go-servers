package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Raghart/go-servers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, req *http.Request) {
	acessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("the request doesn't have an acess token: %v", err))
		return
	}

	chipIdStr := req.PathValue("chirpID")

	pathChirpID, err := uuid.Parse(chipIdStr)
	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("the provided ID is not in a valid uuid format: %v", err))
		return
	}

	databaseChirp, err := cfg.db.GetSingleChirp(context.Background(), pathChirpID)
	if err != nil {
		jsonResponseError(w, http.StatusNotFound, fmt.Sprintf("chip not found: %v", err))
		return
	}

	userID, err := auth.ValidateJWT(acessToken, cfg.secretString)
	if err != nil {
		jsonResponseError(w, 500, fmt.Sprintf("there was a problem while validating the jwt: %v", err))
		return
	}

	if userID != databaseChirp.UserID {
		jsonResponseError(w, 403, "Only the user who created the chirp can delete it")
		return
	}

	err = cfg.db.DeleteChirpByID(context.Background(), databaseChirp.ID)
	if err != nil {
		jsonResponseError(w, 500, fmt.Sprintf("there was a problem while trying to delete the user: %v", err))
		return
	}

	jsonResponse(w, 204, []byte("OK"))
}
