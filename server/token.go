package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

		response, err := s.verifyTokenWithIdP(token)
		if err != nil || response.Active == "false" {
			log.Print("Invalid token ", err.Error())
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

		if !containsScope(response.Scope, routeGroup) {
			log.Print("Insufficient scopes")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) verifyTokenWithIdP(token string) (VerificationResponse, error) {
	req, err := http.NewRequest("POST", "http://localhost:8080/oauth2/token/verify", strings.NewReader(fmt.Sprintf(`{"token": "%s"}`, token)))
	if err != nil {
		return VerificationResponse{}, fmt.Errorf("failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return VerificationResponse{}, fmt.Errorf("failed to verify token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response VerificationErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return VerificationResponse{}, fmt.Errorf("failed to decode response")
		}

		return VerificationResponse{}, fmt.Errorf(response.ErrorDescription)
	}

	var response VerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return VerificationResponse{}, fmt.Errorf("failed to decode response")
	}

	return response, nil
}

func containsScope(tokenScope []string, routeGroup string) bool {
	for _, s := range tokenScope {
		if strings.Contains(s, routeGroup) {
			return true
		}
	}

	return false
}
