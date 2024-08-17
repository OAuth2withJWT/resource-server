package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var publicKey *rsa.PublicKey

func init() {
	pubKeyData, err := os.ReadFile("keys/public_key.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to load public key: %v", err))
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("Failed to decode PEM block containing public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse public key: %v", err))
	}

	publicKey = pubKeyInterface.(*rsa.PublicKey)
}

type VerificationRequest struct {
	Token string `json:"token"`
}

type VerificationErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type VerificationResponse struct {
	Active string   `json:"active"`
	Scope  []string `json:"scope,omitempty"`
}

func (s *Server) protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Missing or invalid Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		parsedToken, err := s.VerifyToken(token)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		currentRoute := mux.CurrentRoute(r)
		routeGroup := ""

		if currentRoute != nil {
			if path, err := currentRoute.GetPathTemplate(); err == nil {
				segments := strings.Split(path, "/")
				routeGroup = segments[2]
			}
		}

		claims, _ := parsedToken.Claims.(jwt.MapClaims)

		scope := claims["scope"].([]interface{})

		if !containsScope(scope, routeGroup) {
			log.Print("Insufficient scopes")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func containsScope(tokenScope []interface{}, routeGroup string) bool {
	for _, s := range tokenScope {
		if scope, ok := s.(string); ok && strings.Contains(scope, routeGroup) {
			return true
		}
	}
	return false
}
