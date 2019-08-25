package serverless

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomBytes creates a slice of random bytes.
func GenerateRandomBytes(length int) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())

	value := make([]byte, length)

	len, err := rand.Read(value)
	if err != nil {
		return []byte{}, err
	}
	if len != length {
		return []byte{}, fmt.Errorf("failed to generate the correct number of bytes")
	}

	return value, nil
}

// GenerateRandomByteString generates a random hex string of bytes for the length given.
func GenerateRandomByteString(length int) (string, error) {
	random, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02x", random), nil
}
