package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Crypto struct{}

var deliminator = "||"

// Generate hash for
func (c *Crypto) Generate(s string) (string, error) {
	salt := uuid.New().String()

	saltedBytes := []byte(s + salt)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash + deliminator + salt, nil
}

func (c *Crypto) Compare(hash string, s string) error {
	// compare password
	parts := strings.Split(hash, deliminator)

	if len(parts) != 2 {
		return errors.New("Invalid hash")
	}
	incoming := []byte(s + parts[1])
	existing := []byte(parts[0])
	return bcrypt.CompareHashAndPassword(existing, incoming)

}
