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

type requestEmail struct {
	Email string `json:"email"`
}

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}
