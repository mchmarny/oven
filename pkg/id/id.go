package id

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

const (
	recordIDPrefix = "id"
)

// NewID generates new ID using UUID v4
func NewID() string {
	return fmt.Sprintf("%s-%s", recordIDPrefix, uuid.New().String())
}

// ToID converts string value into predictable length ID using SHA256 and prefixes it with "id-"
func ToID(v string) string {
	return fmt.Sprintf("%s-%x", recordIDPrefix, sha256.Sum256([]byte(v)))
}
