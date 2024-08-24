package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

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

func (s *Server) protected(next http.Handler, requiredScopes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Missing or invalid Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		err := s.ValidateToken(token, requiredScopes)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) ValidateToken(tokenString string, requiredScopes []string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.RSAConfig.PublicKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to extract claims from token")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != s.RSAConfig.ResourceServer {
		return fmt.Errorf("invalid audience")
	}

	scope, ok := claims["scope"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid or missing scope in token")
	}

	if !containsScopes(scope, requiredScopes) {
		return fmt.Errorf("insufficient scopes")
	}

	return nil
}

func containsScopes(tokenScopes []interface{}, requiredScopes []string) bool {
	scopeSet := make(map[string]struct{})
	for _, s := range tokenScopes {
		if scope, ok := s.(string); ok {
			scopeSet[scope] = struct{}{}
		}
	}

	for _, requiredScope := range requiredScopes {
		if _, exists := scopeSet[requiredScope]; !exists {
			return false
		}
	}

	return true
}
