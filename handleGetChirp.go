package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpIDStr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)

	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("the id provided is not a valid uuid format:'%v'", chirpIDStr))
		return
	}

	dbChirp, err := cfg.db.GetSingleChirp(context.Background(), chirpID)
	if err != nil {
		jsonResponseError(w, 404, fmt.Sprintf("the chirp with the id:'%v' wasn't found: %v", chirpID, err))
		return
	}

	formattedChirp := userChip{
		Id:         dbChirp.ID,
		Created_at: dbChirp.CreatedAt,
		Updated_at: dbChirp.UpdatedAt,
		Body:       dbChirp.Body,
		User_id:    dbChirp.UserID,
	}

	jsonResponse(w, 200, formattedChirp)
}
