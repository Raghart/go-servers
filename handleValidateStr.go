package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Raghart/go-servers/internal/database"
)

func (cfg *apiConfig) handleValidateString(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	strRequest := requestChips{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&strRequest); err != nil {
		log.Printf("error while trying to decode the request: %v", err)
		w.WriteHeader(500)
		return
	}

	if len(strRequest.Body) > 140 {
		jsonResponseError(w, http.StatusBadRequest, "Chirp is too long")
	}

	parsedStr := parseString(strRequest.Body)

	chirpParams := database.CreateChipParams{
		Body:   parsedStr,
		UserID: strRequest.User_id,
	}

	newChirp, err := cfg.db.CreateChip(context.Background(), chirpParams)

	if err != nil {
		jsonResponseError(w, 403, fmt.Sprintf("there was an error while trying to create the new chirp: %v", err))
	}

	newUserChip := userChip{
		Id:         newChirp.ID,
		Created_at: newChirp.CreatedAt,
		Updated_at: newChirp.UpdatedAt,
		Body:       parsedStr,
		User_id:    strRequest.User_id,
	}

	jsonResponse(w, 201, newUserChip)
}

func parseString(str string) string {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	strSlice := strings.Split(str, " ")
	for i, word := range strSlice {
		for _, wordForbid := range forbiddenWords {
			if loweredWord := strings.ToLower(word); loweredWord == wordForbid {
				strSlice[i] = "****"
			}
		}
	}
	return strings.Join(strSlice, " ")
}
