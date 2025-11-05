package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	argon2id.CreateHash()
	return "", nil
}

func CheckPasswordHash(password, hashString string) (bool, error) {
	return true, nil
}
