package generate

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) error {

	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid encoded hash format ")
	}

	saltBase64 := parts[0]
	hashedPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return fmt.Errorf("internal server error")
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("internal server error")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hash) != len(hashedPassword) {
		return fmt.Errorf("incorrect password")
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		return nil
	}

	return fmt.Errorf("incorrect password")
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("Password is blank, Please enter password")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}
