package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/utils"
)

type Middleware interface {
	EnableCORS(h http.Handler) http.Handler
	RequireApiKey(h http.Handler) http.Handler
	RequireAuth(validRoles ...int) func(h http.Handler) http.Handler
}

type middleware struct {
	authService services.AuthService
}

func NewMiddleware(authService services.AuthService) Middleware {
	return &middleware{
		authService,
	}
}

func (m *middleware) EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !utils.IsProdDev() {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			//w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_HOST"))
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CRSF-Token, Authorization, x-api-key")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (m *middleware) RequireApiKey(h http.Handler) http.Handler {
	apiKeyHeader := os.Getenv("APIKeyHeader")
	apiKey := os.Getenv("APIKey")

	decodedApiKey, err := hex.DecodeString(apiKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.URL.Path, "api/v1") {
			h.ServeHTTP(w, r)
			return
		}

		if err != nil {
			fmt.Printf("error occur when decoding api key error : %+v\n", err)
			_ = api.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusInternalServerError)
			return
		}

		apiKeyFromReq, err := m.authService.BearerToken(r, apiKeyHeader)
		if err != nil {
			fmt.Printf("request failed API key authentication error : %+v\n", err)
			_ = api.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)

			return
		}

		if !apiKeyIsValid(apiKeyFromReq, decodedApiKey) {
			hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Printf("failed to parse remote address, error : %+v\n", err)
				hostIP = r.RemoteAddr
			}

			log.Println("no matching API key found", "remoteIP", hostIP)
			_ = api.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)

			return
		}

		h.ServeHTTP(w, r)
	})

}

func (m *middleware) RequireAuth(validRoles ...int) func(h http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		return api.MakeHTTPHandler(func(w http.ResponseWriter, r *http.Request) error {

			jwt, err := m.authService.BearerToken(r, "Authorization")

			if err != nil {
				return api.Error{Err: err.Error(), Status: http.StatusUnauthorized}
			}

			err = m.authService.IsJWTTokenValid(jwt, validRoles...)

			if err != nil {
				return api.Error{Err: err.Error(), Status: http.StatusUnauthorized}

			} else {
				h.ServeHTTP(w, r)
			}

			return nil
		})
	}

}

func apiKeyIsValid(rawKey string, expectedApiKey []byte) bool {
	hash := sha256.Sum256([]byte(rawKey))
	key := string(hash[:])

	contentEqual := subtle.ConstantTimeCompare(expectedApiKey, []byte(key)) == 1

	if contentEqual {
		return true
	}

	return false
}
