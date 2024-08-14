package inmemory

import (
	"github.com/google/uuid"
	"github.com/sosshik/users-service/internal/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	storage := NewInMemory()

	// Setup: Create initial users to test uniqueness constraints
	_, err := storage.CreateUser(models.User{
		Nickname:  "existinguser",
		Email:     "existinguser@example.com",
		FirstName: "Existing",
		LastName:  "User",
		Country:   "Country",
	})
	if err != nil {
		t.Fatalf("Failed to create initial user: %v", err)
	}

	tests := []struct {
		name      string
		input     models.User
		expectErr bool
	}{
		{
			name: "Create valid user",
			input: models.User{
				Nickname:  "newuser",
				Email:     "newuser@example.com",
				FirstName: "New",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: false,
		},
		{
			name: "Create user with existing nickname",
			input: models.User{
				Nickname:  "existinguser", // Existing nickname
				Email:     "different@example.com",
				FirstName: "Another",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: true,
		},
		{
			name: "Create user with existing email",
			input: models.User{
				Nickname:  "differentuser",
				Email:     "existinguser@example.com", // Existing email
				FirstName: "Different",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.CreateUser(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateUser() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	storage := NewInMemory()

	// Create an initial user
	user, _ := storage.CreateUser(models.User{
		Nickname:  "initialuser",
		Email:     "initialuser@example.com",
		FirstName: "Initial",
		LastName:  "User",
		Country:   "Country",
	})

	tests := []struct {
		name      string
		input     models.User
		expectErr bool
	}{
		{
			name: "Update user successfully",
			input: models.User{
				ID:        user.ID,
				Nickname:  "updateduser",
				Email:     "updateduser@example.com",
				FirstName: "Updated",
				LastName:  "User",
				Country:   "UpdatedCountry",
			},
			expectErr: false,
		},
		{
			name: "Update user with existing nickname",
			input: models.User{
				ID:        user.ID,
				Nickname:  "updateduser",
				Email:     "newemail@example.com",
				FirstName: "Another",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: true, // Ensure there's another user with the same nickname
		},
		{
			name: "Update user with existing email",
			input: models.User{
				ID:        user.ID,
				Nickname:  "differentuser",
				Email:     "updateduser@example.com",
				FirstName: "Another",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: true, // Ensure there's another user with the same email
		},
		{
			name: "Update non-existent user",
			input: models.User{
				ID:        uuid.New(), // Invalid ID
				Nickname:  "nonexistent",
				Email:     "nonexistent@example.com",
				FirstName: "Nonexistent",
				LastName:  "User",
				Country:   "Country",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.UpdateUser(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("UpdateUser() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	storage := NewInMemory()

	// Create a user
	user, _ := storage.CreateUser(models.User{
		Nickname:  "getuser",
		Email:     "getuser@example.com",
		FirstName: "Get",
		LastName:  "User",
		Country:   "Country",
	})

	tests := []struct {
		name      string
		input     uuid.UUID
		expectErr bool
	}{
		{
			name:      "Get existing user",
			input:     user.ID,
			expectErr: false,
		},
		{
			name:      "Get non-existent user",
			input:     uuid.New(), // Invalid ID
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.GetUser(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetUser() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	storage := NewInMemory()

	// Create a user
	user, _ := storage.CreateUser(models.User{
		Nickname:  "deleteuser",
		Email:     "deleteuser@example.com",
		FirstName: "Delete",
		LastName:  "User",
		Country:   "Country",
	})

	tests := []struct {
		name      string
		input     uuid.UUID
		expectErr bool
	}{
		{
			name:      "Delete existing user",
			input:     user.ID,
			expectErr: false,
		},
		{
			name:      "Delete non-existent user",
			input:     uuid.New(), // Invalid ID
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DeleteUser(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteUser() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestGetFilteredUsers(t *testing.T) {
	storage := NewInMemory()

	// Setup: Create some users for testing
	users := []models.User{
		{
			Nickname:  "alice",
			Email:     "alice@example.com",
			FirstName: "Alice",
			LastName:  "Smith",
			Country:   "Wonderland",
		},
		{
			Nickname:  "bob",
			Email:     "bob@example.com",
			FirstName: "Bob",
			LastName:  "Johnson",
			Country:   "USA",
		},
		{
			Nickname:  "carol",
			Email:     "carol@example.com",
			FirstName: "Carol",
			LastName:  "Williams",
			Country:   "UK",
		},
	}

	for _, user := range users {
		_, err := storage.CreateUser(user)
		if err != nil {
			t.Fatalf("Failed to create user %v: %v", user, err)
		}
	}

	tests := []struct {
		name          string
		field         string
		value         string
		limit         int
		offset        int
		expected      []models.User
		expectedCount int
	}{
		{
			name:          "Filter by nickname",
			field:         "nickname",
			value:         "bob",
			limit:         10,
			offset:        0,
			expected:      []models.User{users[1]},
			expectedCount: 1, // Corrected expected count
		},
		{
			name:          "Filter by email",
			field:         "email",
			value:         "alice@example.com",
			limit:         10,
			offset:        0,
			expected:      []models.User{users[0]},
			expectedCount: 1, // Corrected expected count
		},
		{
			name:          "Filter by first name",
			field:         "first_name",
			value:         "Bob",
			limit:         10,
			offset:        0,
			expected:      []models.User{users[1]},
			expectedCount: 1, // Corrected expected count
		},
		{
			name:          "Filter by country",
			field:         "country",
			value:         "USA",
			limit:         1,
			offset:        0,
			expected:      []models.User{users[1]},
			expectedCount: 1, // Corrected expected count
		},
		{
			name:          "Pagination",
			field:         "nickname",
			value:         "",
			limit:         2,
			offset:        1,
			expected:      []models.User{users[1], users[2]},
			expectedCount: 3,
		},
		{
			name:          "No results",
			field:         "nickname",
			value:         "nonexistent",
			limit:         10,
			offset:        0,
			expected:      []models.User{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, count, err := storage.GetFilteredUsers(tt.field, tt.value, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("GetFilteredUsers() error = %v", err)
				return
			}

			if count != tt.expectedCount {
				t.Errorf("GetFilteredUsers() returned count = %d, expected %d", count, tt.expectedCount)
			}

			// Compare results while ignoring UUIDs and timestamps
			for i, user := range tt.expected {
				if i >= len(result) {
					t.Errorf("Expected user %v not found in results", user)
					continue
				}
				if !usersEqualIgnoringDynamicFields(result[i], user) {
					t.Errorf("GetFilteredUsers() user %d = %v, expected %v", i, result[i], user)
				}
			}
		})
	}
}

// Helper function to compare users while ignoring UUID and timestamps
func usersEqualIgnoringDynamicFields(a, b models.User) bool {
	return a.Nickname == b.Nickname &&
		a.Email == b.Email &&
		a.FirstName == b.FirstName &&
		a.LastName == b.LastName &&
		a.Country == b.Country
}
