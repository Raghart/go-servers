package main

import (
	"context"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetChirps(context.Background())

	if err != nil {
		jsonResponseError(w, http.StatusBadRequest, fmt.Sprintf("there was a problem while trying to get the chirps: %v", err))
		return
	}

	var chirpSlice []userChip
	for _, chirp := range dbChirps {
		chirpSlice = append(chirpSlice, userChip{
			Id:         chirp.ID,
			Created_at: chirp.CreatedAt,
			Updated_at: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_id:    chirp.UserID,
		})
	}

	jsonResponse(w, 200, chirpSlice)
}
