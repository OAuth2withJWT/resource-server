package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type RSAConfig struct {
	PublicKey      *rsa.PublicKey
	ResourceServer string
}

func LoadRSAConfig() RSAConfig {
	resourceServer := os.Getenv("RESOURCE_SERVER")
	if resourceServer == "" {
		log.Fatalf("RESOURCE_SERVER environment variable not set")
	}

	config := RSAConfig{
		ResourceServer: resourceServer,
	}
	publicKey, err := LoadPublicKey("keys/public_key.pem")
	if err != nil {
		log.Printf("Failed to load public key, continuing with default config: %v", err)
	} else {
		config.PublicKey = publicKey
	}

	return config
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pubKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}

	return rsaPublicKey, nil
}
