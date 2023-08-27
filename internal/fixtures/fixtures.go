package fixtures

import "github.com/rustoma/octo-pulse/internal/models"

func CreateUser(fn, ln string, admin bool) *models.User {

	user := &models.User{
		ID: 1,
	}

	return user
}
