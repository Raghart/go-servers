package main

import (
	"time"

	"github.com/google/uuid"
)

type cleanedResponse struct {
	Cleaned_body string `json:"cleaned_body"`
}

type requestChips struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

type userChip struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}

type requestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id            uuid.UUID `json:"id"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Is_chirpy_red bool      `json:"is_chirpy_red"`
}

type UserLogin struct {
	Id            uuid.UUID `json:"id"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
	Is_chirpy_red bool      `json:"is_chirpy_red"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type webHookRequest struct {
	Event string `json:"event"`
	Data  struct {
		User_id string `json:"user_id"`
	}
}
