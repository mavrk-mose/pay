package middleware

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
)

// Sign creates a signature for a binary file located at binaryPath using the provided RSA private key.
// It returns the signature as a byte slice and any error encountered.
func Sign(binaryPath string, privateKey *rsa.PrivateKey) ([]byte, error) {
	binaryData, err := os.ReadFile(binaryPath)
	if err != nil {
	 return nil, fmt.Errorf("read binary data: %v", err)
	}
   
	hashed := sha256.Sum256(binaryData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
	 return nil, fmt.Errorf("sign data: %v", err)
	}
	return signature, nil
}

// VerifySignature checks the signature of a binary file located at binaryPath against a signature file at signaturePath.
// It uses the provided RSA public key for verification. It returns any error encountered in the verification process.
func VerifySignature(binaryPath, signaturePath string, publicKey *rsa.PublicKey) error {
	binaryData, err := os.ReadFile(binaryPath)
	if err != nil {
	 return fmt.Errorf("read binary data: %v", err)
	}
   
	signature, err := os.ReadFile(signaturePath)
	if err != nil {
	 return fmt.Errorf("read signature data: %v", err)
	}
   
	hashed := sha256.Sum256(binaryData)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
	 return fmt.Errorf("verify signature: %v", err)
	}
	return nil
}