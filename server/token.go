package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
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

func (s *Server) protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Missing or invalid Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		parsedToken, err := s.ValidateToken(token)
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

func (s *Server) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.RSAConfig.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != s.RSAConfig.ResourceServer {
		return nil, fmt.Errorf("invalid audience")
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
