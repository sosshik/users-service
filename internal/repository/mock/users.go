package mocks

import (
	"github.com/google/uuid"
	"github.com/sosshik/users-service/internal/models"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetUser(id uuid.UUID) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) NicknameOrEmailExists(nickname, email string) (bool, error) {
	args := m.Called(nickname, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetFilteredUsers(field, value string, limit, offset int) ([]models.User, int, error) {
	args := m.Called(field, value, limit, offset)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}
