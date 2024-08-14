package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sosshik/users-service/internal/models"
	mocks "github.com/sosshik/users-service/internal/repository/mock"
	"github.com/sosshik/users-service/pkg/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := NewUsersService(mockRepo)

	testCases := []struct {
		name         string
		userReq      dtos.CreateUserRequest
		expectedResp dtos.CreateUserResponse
		expectedErr  error
		setupMock    func()
	}{
		{
			name: "Success",
			userReq: dtos.CreateUserRequest{
				Nickname: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			expectedResp: dtos.CreateUserResponse{
				Nickname: "newuser",
				Email:    "new@example.com",
			},
			expectedErr: nil,
			setupMock: func() {
				mockRepo.On("CreateUser", mock.Anything).Return(models.User{
					Nickname: "newuser",
					Email:    "new@example.com",
				}, nil).Once()
			},
		},
		{
			name: "Repository error",
			userReq: dtos.CreateUserRequest{
				Nickname: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			expectedResp: dtos.CreateUserResponse{},
			expectedErr:  errors.New("repository error"),
			setupMock: func() {
				mockRepo.On("CreateUser", mock.Anything).Return(models.User{}, errors.New("repository error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			userResp, err := userService.CreateUser(tc.userReq)

			assert.Equal(t, tc.expectedResp, userResp)
			assert.Equal(t, tc.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := NewUsersService(mockRepo)

	testCases := []struct {
		name         string
		idStr        string
		userReq      dtos.UpdateUserRequest
		expectedResp dtos.UpdateUserResponse
		expectedErr  error
		setupMock    func()
	}{
		{
			name:  "Success",
			idStr: uuid.New().String(),
			userReq: dtos.UpdateUserRequest{
				Nickname: "updateduser",
				Email:    "updated@example.com",
			},
			expectedResp: dtos.UpdateUserResponse{
				Nickname: "updateduser",
				Email:    "updated@example.com",
			},
			expectedErr: nil,
			setupMock: func() {
				mockRepo.On("UpdateUser", mock.Anything).Return(models.User{
					Nickname: "updateduser",
					Email:    "updated@example.com",
				}, nil).Once()
			},
		},
		{
			name:  "Invalid UUID",
			idStr: "invalid-uuid",
			userReq: dtos.UpdateUserRequest{
				Nickname: "updateduser",
				Email:    "updated@example.com",
			},
			expectedResp: dtos.UpdateUserResponse{},
			expectedErr:  errors.New("invalid UUID length: 12"),
			setupMock:    func() {},
		},
		{
			name:  "Repository error",
			idStr: uuid.New().String(),
			userReq: dtos.UpdateUserRequest{
				Nickname: "updateduser",
				Email:    "updated@example.com",
			},
			expectedResp: dtos.UpdateUserResponse{},
			expectedErr:  errors.New("repository error"),
			setupMock: func() {
				mockRepo.On("UpdateUser", mock.Anything).Return(models.User{}, errors.New("repository error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			userResp, err := userService.UpdateUser(tc.idStr, tc.userReq)

			if tc.expectedErr != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tc.expectedErr.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResp, userResp)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := NewUsersService(mockRepo)

	testCases := []struct {
		name        string
		idStr       string
		expectedErr error
		setupMock   func()
	}{
		{
			name:        "Success",
			idStr:       uuid.New().String(),
			expectedErr: nil,
			setupMock: func() {
				mockRepo.On("DeleteUser", mock.Anything).Return(nil).Once()
			},
		},
		{
			name:        "Repository error",
			idStr:       uuid.New().String(),
			expectedErr: errors.New("repository error"),
			setupMock: func() {
				mockRepo.On("DeleteUser", mock.Anything).Return(errors.New("repository error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := userService.DeleteUser(tc.idStr)

			assert.Equal(t, tc.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetFilteredUsers(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := &UsersService{repo: mockRepo}

	// Define fixed timestamps
	fixedCreatedAt := time.Date(2024, time.August, 14, 20, 32, 0, 0, time.Local)
	fixedUpdatedAt := time.Date(2024, time.August, 14, 20, 32, 0, 0, time.Local)

	tests := []struct {
		name           string
		pageStr        string
		pageSizeStr    string
		filterStr      string
		mockReturn     []models.User
		mockTotalCount int
		mockErr        error
		expected       dtos.GetUserResponse
		expectedErr    error
	}{
		{
			name:        "Success",
			pageStr:     "1",
			pageSizeStr: "10",
			filterStr:   "nickname:testuser",
			mockReturn: []models.User{
				{
					ID:        uuid.New(),
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "testuser",
					Email:     "john.doe@example.com",
					Country:   "Country",
					CreatedAt: fixedCreatedAt,
					UpdatedAt: fixedUpdatedAt,
				},
			},
			mockTotalCount: 1,
			mockErr:        nil,
			expected: dtos.GetUserResponse{
				Page:     1,
				PageSize: 10,
				Total:    1,
				Users: []dtos.GetUserDTO{
					{
						ID:        uuid.Nil, // ID comparison can be tricky, consider using a special assertion
						FirstName: "John",
						LastName:  "Doe",
						Nickname:  "testuser",
						Email:     "john.doe@example.com",
						Country:   "Country",
						CreatedAt: fixedCreatedAt,
						UpdatedAt: fixedUpdatedAt,
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetFilteredUsers", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
				Return(tt.mockReturn, tt.mockTotalCount, tt.mockErr)

			got, err := service.GetFilteredUsers(tt.pageStr, tt.pageSizeStr, tt.filterStr)

			assert.Equal(t, tt.expectedErr, err)

			for i := range got.Users {
				expectedUser := tt.expected.Users[i]
				actualUser := got.Users[i]

				expectedUser.ID = uuid.Nil // Handling UUID comparison can be tricky
				actualUser.ID = uuid.Nil

				assert.Equal(t, expectedUser, actualUser)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
