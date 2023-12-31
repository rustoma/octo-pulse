package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rustoma/octo-pulse/internal/dto"
	e "github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(*dto.AuthLogin) (*dto.AuthUser, *http.Cookie, error)
	Logout(*dto.LogoutRequest) (*http.Cookie, error)
	RefreshToken(refreshTokenRequest *dto.RefreshTokenRequest) (string, error)
	CheckPassword(password string, hashedPassword string) error
	HashPassword(password string) (string, error)
	BearerToken(r *http.Request, header string) (string, error)
	IsJWTTokenValid(tokenString string, validRoles ...int) error
	validateUserRole(userRoles int, validRoles []int) error
	parseToken(jwtString string) (*jwt.Token, error)
	generateJWTToken(claims JWTClaims) (string, error)
	createTokenExpirationTimeForJWTRefreshToken() *jwt.NumericDate
	createTokenExpirationTimeForJWTToken() *jwt.NumericDate
}

type authService struct {
	userStore storage.UserStore
	userRoles models.UserRoles
}

type JWTClaims struct {
	Email string `json:"email"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userStore storage.UserStore) AuthService {
	return &authService{userStore: userStore, userRoles: models.UserRoles{
		Admin:  1,
		Editor: 2,
	}}
}

func (a *authService) Login(userCredentials *dto.AuthLogin) (*dto.AuthUser, *http.Cookie, error) {
	user, err := a.userStore.GetUserByEmail(userCredentials.Email)

	if err != nil {
		return nil, nil, e.NotFound{Err: "user not found"}
	}

	if !user.IsEnabled {
		return nil, nil, e.BadRequest{Err: "user not found"}
	}

	err = a.CheckPassword(userCredentials.Password, user.PasswordHash)

	if err != nil {
		return nil, nil, e.BadRequest{Err: "bad user password"}
	}

	JWTTokenClaims := JWTClaims{
		Email: userCredentials.Email,
		Role:  user.RoleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: a.createTokenExpirationTimeForJWTToken(),
			Issuer:    os.Getenv("SERVER_IP"),
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
	}

	refreshTokenClaims := JWTClaims{
		Email: userCredentials.Email,
		Role:  user.RoleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: a.createTokenExpirationTimeForJWTRefreshToken(),
			Issuer:    os.Getenv("SERVER_IP"),
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
	}

	encodedJWT, _ := a.generateJWTToken(JWTTokenClaims)
	encodedRefreshToken, _ := a.generateJWTToken(refreshTokenClaims)

	_, err = a.userStore.UpdateRefreshToken(user.ID, encodedRefreshToken)

	if err != nil {
		return nil, nil, err
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    encodedRefreshToken,
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		Secure:   utils.IsProdDev(),
		Path:     "/",
		SameSite: 4,
	}

	user.RefreshToken = encodedRefreshToken

	return &dto.AuthUser{User: user, AccessToken: encodedJWT}, &cookie, nil
}

func (a *authService) Logout(logoutRequest *dto.LogoutRequest) (*http.Cookie, error) {

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   utils.IsProdDev(),
		Path:     "/",
		SameSite: 4,
	}

	user, err := a.userStore.SelectUserByRefreshToken(logoutRequest.RefreshToken)

	if err != nil {
		return cookie, e.NotFound{Err: "user not found"}
	}

	_, err = a.userStore.UpdateRefreshToken(user.ID, "")

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (a *authService) RefreshToken(refreshTokenRequest *dto.RefreshTokenRequest) (string, error) {

	user, err := a.userStore.SelectUserByRefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		return "", e.NotFound{Err: err.Error()}
	}

	if !user.IsEnabled {
		return "", e.NotFound{Err: err.Error()}
	}

	token, err := a.parseToken(refreshTokenRequest.RefreshToken)

	if err != nil {
		return "", e.Unauthorized{Err: err.Error()}
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		userEmail := claims.Email

		if user.Email != userEmail {
			return "", e.Unauthorized{Err: err.Error()}
		}

		JWTTokenClaims := JWTClaims{
			Email: user.Email,
			Role:  user.RoleId,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: a.createTokenExpirationTimeForJWTToken(),
				Issuer:    os.Getenv("SERVER_IP"),
				IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
			},
		}

		encodedJWT, _ := a.generateJWTToken(JWTTokenClaims)
		return encodedJWT, nil
	} else {
		return "", e.Unauthorized{Err: "JWT Claims are not correct"}
	}
}

func (claims JWTClaims) Validate() error {
	if claims.Email == "" {
		return errors.New("user name claims are missing")
	}
	return nil
}

func (a *authService) createTokenExpirationTimeForJWTToken() *jwt.NumericDate {
	ttl := 30 * time.Minute
	expirationTime := time.Now().UTC().Add(ttl)
	return &jwt.NumericDate{Time: expirationTime}
}

func (a *authService) createTokenExpirationTimeForJWTRefreshToken() *jwt.NumericDate {
	ttl := 24 * time.Hour
	expirationTime := time.Now().UTC().Add(ttl)
	return &jwt.NumericDate{Time: expirationTime}
}

func (a *authService) generateJWTToken(claims JWTClaims) (string, error) {
	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte(os.Getenv("JWT_SECRET"))

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(key)

	if err != nil {
		return "", err
	}

	return s, nil
}

func (a *authService) parseToken(jwtString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(jwtString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func (a *authService) validateUserRole(userRole int, validRoles []int) error {
	if userRole == a.userRoles.Admin {
		return nil
	}

	if len(validRoles) > 0 {
		for _, role := range validRoles {
			if role == userRole {
				return nil
			}
		}
		return errors.New("you do not have enough permissions")
	}

	return nil
}

func (a *authService) IsJWTTokenValid(tokenString string, validRoles ...int) error {

	var err error

	token, err := a.parseToken(tokenString)

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		userRole := claims.Role
		return a.validateUserRole(userRole, validRoles)
	} else {
		return errors.New("JWT Claims are not correct")
	}
}

func (a *authService) BearerToken(r *http.Request, header string) (string, error) {
	rawToken := r.Header.Get(header)
	pieces := strings.SplitN(rawToken, " ", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}

func (a *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (a *authService) CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
