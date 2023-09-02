package dto

import "github.com/rustoma/octo-pulse/internal/models"

type AuthLogin struct {
	Email    string
	Password string
}

type AuthUser struct {
	User        *models.User `json:"user"`
	AccessToken string       `json:"accessToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
