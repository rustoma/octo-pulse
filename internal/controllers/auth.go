package controllers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/services"
)

type JWTClaims struct {
	UserName string `json:"user_name"`
	Roles    []int  `json:"roles"`
	jwt.RegisteredClaims
}

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var userCredentials *dto.AuthLogin

	err := api.ReadJSON(w, r, &userCredentials)

	if err != nil {
		return api.Error{Err: "bad login request", Status: http.StatusBadRequest}
	}

	authUser, cookie, err := c.authService.Login(userCredentials)

	if err != nil {
		return api.Error{Err: "Cannot login", Status: api.HandleErrorStatus(err)}
	}

	http.SetCookie(w, cookie)

	return api.WriteJSON(w, http.StatusOK, authUser)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	var logoutRequest *dto.LogoutRequest

	err := api.ReadJSON(w, r, &logoutRequest)

	if err != nil {
		return api.WriteJSON(w, http.StatusNoContent, "")
	}

	cookie, err := c.authService.Logout(logoutRequest)

	if err != nil {
		if cookie != nil {
			http.SetCookie(w, cookie)
		}
		return err
	}

	http.SetCookie(w, cookie)
	return api.WriteJSON(w, http.StatusNoContent, "Logout successful")
}

func (c *AuthController) HandleRefreshToken(w http.ResponseWriter, r *http.Request) error {
	var refreshTokenRequest *dto.RefreshTokenRequest

	err := api.ReadJSON(w, r, &refreshTokenRequest)

	if err != nil {
		return api.Error{Err: "refresh token not found", Status: api.HandleErrorStatus(err)}
	}

	encodedJWT, err := c.authService.RefreshToken(refreshTokenRequest)

	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, dto.RefreshTokenResponse{AccessToken: encodedJWT})

}
