package inmemory

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sosshik/users-service/internal/models"
	"strings"
	"sync"
	"time"
)

type InMemoryStorage struct {
	mu            sync.RWMutex
	users         []models.User
	idIndex       map[uuid.UUID]*models.User
	nicknameIndex map[string]*models.User
	emailIndex    map[string]*models.User
}

// NewInMemory creates a new instance of InMemoryStorage with initialized data structures
func NewInMemory() *InMemoryStorage {
	return &InMemoryStorage{
		users:         make([]models.User, 0),
		idIndex:       make(map[uuid.UUID]*models.User),
		nicknameIndex: make(map[string]*models.User),
		emailIndex:    make(map[string]*models.User),
	}
}

// CreateUser adds a new user to the in-memory storage
func (s *InMemoryStorage) CreateUser(user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if a user with the given nickname or email already exists
	exists, err := s.nicknameOrEmailExists(user.Nickname, user.Email)
	if err != nil || exists {
		return models.User{}, err
	}

	// Assign a new UUID and set timestamps
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Append the new user to the list and update indexes
	s.users = append(s.users, user)
	s.idIndex[user.ID] = &s.users[len(s.users)-1]
	s.nicknameIndex[user.Nickname] = &s.users[len(s.users)-1]
	s.emailIndex[user.Email] = &s.users[len(s.users)-1]

	return user, nil
}

// NicknameOrEmailExists checks if a user with the given nickname or email already exists
func (s *InMemoryStorage) NicknameOrEmailExists(nickname, email string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nicknameOrEmailExists(nickname, email)
}

// nicknameOrEmailExists is a helper function that checks existence of a user by nickname or email
func (s *InMemoryStorage) nicknameOrEmailExists(nickname, email string) (bool, error) {
	if _, exists := s.nicknameIndex[nickname]; exists {
		return true, errors.New("user with this username already exists")
	}

	if _, exists := s.emailIndex[email]; exists {
		return true, errors.New("user with this email already exists")
	}

	return false, nil
}

// GetUser retrieves a user by their ID
func (s *InMemoryStorage) GetUser(id uuid.UUID) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, found := s.idIndex[id]
	if !found {
		return models.User{}, errors.New("user not found")
	}

	return *user, nil
}

// UpdateUser modifies an existing user's details
func (s *InMemoryStorage) UpdateUser(user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user exists
	if _, found := s.idIndex[user.ID]; found {
		// Validate that new nickname/email does not exist
		if exists, err := s.nicknameOrEmailExists(user.Nickname, user.Email); err != nil || exists {
			return models.User{}, err
		}
		// Update timestamps and copy data
		user.UpdatedAt = time.Now()
		oldUser := *s.idIndex[user.ID]
		err := copier.CopyWithOption(s.idIndex[user.ID], &user, copier.Option{IgnoreEmpty: true})
		if err != nil {
			return models.User{}, err
		}
		// Update indexes if nickname/email changed
		if user.Nickname != oldUser.Nickname {
			delete(s.nicknameIndex, oldUser.Nickname)
			s.nicknameIndex[user.Nickname] = s.idIndex[user.ID]
		}
		if user.Email != oldUser.Email {
			delete(s.emailIndex, oldUser.Email)
			s.emailIndex[user.Email] = s.idIndex[user.ID]
		}
		return *s.idIndex[user.ID], nil
	}

	return models.User{}, errors.New("user not found, unable to update")
}

// DeleteUser removes a user from storage by their ID
func (s *InMemoryStorage) DeleteUser(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, found := s.idIndex[id]
	if !found {
		return errors.New("user not found, unable to delete")
	}

	// Remove user from the list and indexes
	for i, u := range s.users {
		if u.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			break
		}
	}

	delete(s.idIndex, id)
	delete(s.nicknameIndex, user.Nickname)
	delete(s.emailIndex, user.Email)

	return nil
}

// GetFilteredUsers retrieves users based on a filter and pagination parameters
func (s *InMemoryStorage) GetFilteredUsers(field, value string, limit, offset int) ([]models.User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.User

	// Normalize the value to lowercase
	value = strings.ToLower(value)

	// Filter users based on the provided field and value
	for _, user := range s.users {
		if needToIncludeUser(user, field, value) {
			result = append(result, user)
		}
	}

	// Implement pagination
	start := offset
	if start > len(result) {
		start = len(result)
	}

	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], len(result), nil
}

// needToIncludeUser checks if a user should be included in the result based on the filter criteria
func needToIncludeUser(user models.User, field, value string) bool {
	if field == "" || value == "" {
		return true
	}
	switch field {
	case "nickname":
		return strings.Contains(strings.ToLower(user.Nickname), value)
	case "email":
		return strings.Contains(strings.ToLower(user.Email), value)
	case "first_name":
		return strings.Contains(strings.ToLower(user.FirstName), value)
	case "last_name":
		return strings.Contains(strings.ToLower(user.LastName), value)
	case "country":
		return strings.Contains(strings.ToLower(user.Country), value)
	}
	return false
}
