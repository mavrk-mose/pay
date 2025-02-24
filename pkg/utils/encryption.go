package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateChecksum(transactionID string, amount float64, currency string, timestamp string) string {
	data := fmt.Sprintf("%s|%.2f|%s|%s", transactionID, amount, currency, timestamp)
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}