package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("error while parsing the payload: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
	return nil
}

func jsonResponseError(w http.ResponseWriter, statusCode int, msg string) error {
	return jsonResponse(w, statusCode, map[string]string{"error": msg})
}
