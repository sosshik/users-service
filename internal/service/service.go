package service

import (
	"github.com/sosshik/users-service/internal/dtos"
	"github.com/sosshik/users-service/internal/repository"
)

type Users interface {
	CreateUser(userReq dtos.CreateUserRequest) (dtos.CreateUserResponse, error)
	UpdateUser(id string, userReq dtos.UpdateUserRequest) (dtos.UpdateUserResponse, error)
	DeleteUser(idStr string) error
	GetFilteredUsers(pageStr, pageSizeStr, filterStr string) (dtos.GetUserResponse, error)
}

type Service struct {
	Users
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Users: NewUsersService(repo),
	}
}
