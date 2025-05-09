package middleware

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	signatureHeader = "Signature"
	filePath        = "certs/public.pem"
)

func SignatureMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		signatureHeader := c.GetHeader("signature")
		if signatureHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing signature header"})
			c.Abort()
			return
		}

		signature, err := base64.StdEncoding.DecodeString(signatureHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature encoding"})
			c.Abort()
			return
		}

		if err := VerifySignature(body, signature, publicKey); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Sign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("sign data: %v", err)
	}
	return signature, nil
}

func VerifySignature(data, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return fmt.Errorf("verify signature: %v", err)
	}
	return nil
}

func LoadPublicKey() (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}
