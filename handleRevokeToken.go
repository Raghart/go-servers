package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Raghart/go-servers/internal/auth"
)

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("there is no valid request token: %v", err))
		return
	}

	_, err = cfg.db.RevokeToken(context.Background(), refreshToken)

	if err != nil {
		jsonResponseError(w, 401, fmt.Sprintf("there was an error while trying to revoke the token: %v", err))
		return
	}

	jsonResponse(w, 204, []byte("OK"))
}
