package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateHash() string {
	timestamp := time.Now().Unix()

	randBytes := make([]byte,4)
	rand.Read(randBytes)

	prefix := fmt.Sprintf("%d_%s", timestamp, hex.EncodeToString(randBytes))

	return prefix
}