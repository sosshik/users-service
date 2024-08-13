package repository

import (
	"github.com/google/uuid"
	"github.com/sosshik/users-service/internal/models"
	"github.com/sosshik/users-service/internal/repository/inmemory"
)

type Users interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(id uuid.UUID) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(id uuid.UUID) error
	NicknameOrEmailExists(nickname, email string) (bool, error)
	GetFilteredUsers(field, value string, limit, offset int) ([]models.User, int, error)
}

type Repository struct {
	Users
}

func NewRepository() *Repository {
	return &Repository{
		Users: inmemory.NewInMemory(),
	}
}
