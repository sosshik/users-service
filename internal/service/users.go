package service

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sosshik/users-service/internal/models"
	"github.com/sosshik/users-service/internal/repository"
	"github.com/sosshik/users-service/pkg/dtos"
	"github.com/sosshik/users-service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UsersService struct {
	repo repository.Users
}

// NewUsersService creates a new instance of UsersService with the given repository
func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{repo}
}

// CreateUser processes the request to create a new user
func (u *UsersService) CreateUser(userReq dtos.CreateUserRequest) (dtos.CreateUserResponse, error) {
	var userResp dtos.CreateUserResponse
	var user models.User

	// Copy data from request DTO to model
	err := copier.Copy(&user, &userReq)
	if err != nil {
		return userResp, err
	}

	// Hash the user's password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return userResp, err
	}
	user.Password = string(hash)

	// Create the user in the repository
	user, err = u.repo.CreateUser(user)
	if err != nil {
		return userResp, err
	}

	// Copy the created user data to response DTO
	err = copier.Copy(&userResp, &user)

	return userResp, err
}

// UpdateUser processes the request to update an existing user
func (u *UsersService) UpdateUser(idStr string, userReq dtos.UpdateUserRequest) (dtos.UpdateUserResponse, error) {
	// Parse user ID from string
	id, err := uuid.Parse(idStr)
	if err != nil {
		return dtos.UpdateUserResponse{}, err
	}

	var userResp dtos.UpdateUserResponse
	var user models.User

	// Copy data from request DTO to model
	err = copier.Copy(&user, &userReq)
	if err != nil {
		return userResp, err
	}
	user.ID = id

	// Update the user in the repository
	user, err = u.repo.UpdateUser(user)
	if err != nil {
		return userResp, err
	}

	// Copy the updated user data to response DTO
	err = copier.Copy(&userResp, &user)

	return userResp, err
}

// DeleteUser processes the request to delete a user by ID
func (u *UsersService) DeleteUser(idStr string) error {
	// Parse user ID from string
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	// Delete the user from the repository
	return u.repo.DeleteUser(id)
}

// GetFilteredUsers retrieves users based on filter and pagination parameters
func (u *UsersService) GetFilteredUsers(pageStr, pageSizeStr, filterStr string) (dtos.GetUserResponse, error) {
	// Convert page number from string to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return dtos.GetUserResponse{}, err
	}
	if page < 1 {
		page = 1
	}

	// Convert page size from string to integer
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return dtos.GetUserResponse{}, err
	}
	if pageSize < 10 {
		pageSize = 10
	}

	// Process filter to get field and value
	field, value := utils.ProcessFilter(filterStr)

	// Retrieve filtered users from the repository
	users, totalFilteredUsers, err := u.repo.GetFilteredUsers(field, value, pageSize, pageSize*(page-1))
	if err != nil {
		return dtos.GetUserResponse{}, err
	}

	var userDTOs []dtos.GetUserDTO
	// Copy user data from model to DTOs
	err = copier.Copy(&userDTOs, users)
	if err != nil {
		return dtos.GetUserResponse{}, err
	}

	// Return the paginated and filtered user data
	return dtos.GetUserResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    totalFilteredUsers,
		Users:    userDTOs,
	}, nil
}
