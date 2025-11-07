package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	authorIdStr := req.URL.Query().Get("author_id")
	sortValue := req.URL.Query().Get("sort")

	if len(sortValue) == 0 {
		sortValue = "asc"
	}

	if sortValue != "asc" && sortValue != "desc" {
		jsonResponseError(w, 400, "Unvalid sort option, only 'asc' or 'desc' are avaible")
		return
	}

	if len(authorIdStr) > 0 {
		authorId, err := uuid.Parse(authorIdStr)

		if err != nil {
			jsonResponseError(w, 400, fmt.Sprintf("the id provided is not in a valid uuid format: %v", err))
			return
		}

		userChirps, err := cfg.db.GetChirpsByUserID(context.Background(), authorId)
		if err != nil {
			jsonResponseError(w, 400, fmt.Sprintf("the provided userID doesn't exist in the database: %v", err))
		}

		var chirpSlice []userChip
		for _, chirp := range userChirps {
			chirpSlice = append(chirpSlice, userChip{
				Id:         chirp.ID,
				Created_at: chirp.CreatedAt,
				Updated_at: chirp.UpdatedAt,
				Body:       chirp.Body,
				User_id:    chirp.UserID,
			})
		}

		sortedSlice := sortChirps(chirpSlice, sortValue)

		jsonResponse(w, http.StatusOK, sortedSlice)

	} else {
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

		sortedSlice := sortChirps(chirpSlice, sortValue)

		jsonResponse(w, http.StatusOK, sortedSlice)
	}
}

func sortChirps(usersChips []userChip, sortType string) []userChip {
	sort.Slice(usersChips, func(i, j int) bool {
		if sortType == "desc" {
			return usersChips[i].Created_at.After(usersChips[j].Created_at)
		}
		return usersChips[i].Created_at.Before(usersChips[j].Created_at)
	})

	return usersChips
}
