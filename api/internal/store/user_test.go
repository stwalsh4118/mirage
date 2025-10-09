package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&User{})
	require.NoError(t, err)

	return db
}

func TestUserCreation(t *testing.T) {
	db := setupTestDB(t)

	// Test creating user with all fields
	user := User{
		ClerkUserID:     "user_test123",
		Email:           "test@example.com",
		FirstName:       stringPtr("John"),
		LastName:        stringPtr("Doe"),
		ProfileImageURL: stringPtr("https://example.com/avatar.jpg"),
		IsActive:        true,
	}

	err := db.Create(&user).Error
	require.NoError(t, err)

	// Verify user was created with correct data
	var retrieved User
	err = db.First(&retrieved, "id = ?", user.ID).Error
	require.NoError(t, err)

	assert.Equal(t, "user_test123", retrieved.ClerkUserID)
	assert.Equal(t, "test@example.com", retrieved.Email)
	assert.Equal(t, "John", *retrieved.FirstName)
	assert.Equal(t, "Doe", *retrieved.LastName)
	assert.Equal(t, "https://example.com/avatar.jpg", *retrieved.ProfileImageURL)
	assert.True(t, retrieved.IsActive)
	assert.NotZero(t, retrieved.CreatedAt)
	assert.NotZero(t, retrieved.UpdatedAt)
}

func TestUserCreationMinimalFields(t *testing.T) {
	db := setupTestDB(t)

	// Test creating user with only required fields
	user := User{
		ClerkUserID: "user_minimal123",
		Email:       "minimal@example.com",
		// FirstName, LastName, ProfileImageURL are nil
		// IsActive should default to true
	}

	err := db.Create(&user).Error
	require.NoError(t, err)

	// Verify defaults and nil handling
	var retrieved User
	err = db.First(&retrieved, "clerk_user_id = ?", user.ClerkUserID).Error
	require.NoError(t, err)

	assert.Equal(t, "user_minimal123", retrieved.ClerkUserID)
	assert.Equal(t, "minimal@example.com", retrieved.Email)
	assert.Nil(t, retrieved.FirstName)
	assert.Nil(t, retrieved.LastName)
	assert.Nil(t, retrieved.ProfileImageURL)
	assert.True(t, retrieved.IsActive) // Should default to true
	assert.NotZero(t, retrieved.CreatedAt)
	assert.NotZero(t, retrieved.UpdatedAt)
}

func TestUserUniqueConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Create first user
	user1 := User{
		ClerkUserID: "user_unique123",
		Email:       "unique@example.com",
	}
	err := db.Create(&user1).Error
	require.NoError(t, err)

	// Try to create user with same ClerkUserID
	user2 := User{
		ClerkUserID: "user_unique123", // Same as user1
		Email:       "different@example.com",
	}
	err = db.Create(&user2).Error
	assert.Error(t, err) // Should fail due to unique constraint

	// Try to create user with same email
	user3 := User{
		ClerkUserID: "user_different456",
		Email:       "unique@example.com", // Same as user1
	}
	err = db.Create(&user3).Error
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestUserFullName(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name: "both names present",
			user: User{
				FirstName: stringPtr("John"),
				LastName:  stringPtr("Doe"),
				Email:     "john@example.com",
			},
			expected: "John Doe",
		},
		{
			name: "only first name",
			user: User{
				FirstName: stringPtr("John"),
				Email:     "john@example.com",
			},
			expected: "John",
		},
		{
			name: "only last name",
			user: User{
				LastName: stringPtr("Doe"),
				Email:    "john@example.com",
			},
			expected: "Doe",
		},
		{
			name: "no names, fallback to email",
			user: User{
				Email: "john@example.com",
			},
			expected: "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.FullName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserLastSeenAt(t *testing.T) {
	db := setupTestDB(t)

	user := User{
		ClerkUserID: "user_seen123",
		Email:       "seen@example.com",
		LastSeenAt:  nil, // Initially nil
	}

	// Create user
	err := db.Create(&user).Error
	require.NoError(t, err)

	// Update LastSeenAt
	now := time.Now()
	err = db.Model(&user).Update("last_seen_at", now).Error
	require.NoError(t, err)

	// Verify update
	var retrieved User
	err = db.First(&retrieved, "clerk_user_id = ?", user.ClerkUserID).Error
	require.NoError(t, err)

	assert.NotNil(t, retrieved.LastSeenAt)
	assert.True(t, retrieved.LastSeenAt.Equal(now))
}

func TestUserSoftDelete(t *testing.T) {
	db := setupTestDB(t)

	user := User{
		ClerkUserID: "user_softdelete123",
		Email:       "softdelete@example.com",
		IsActive:    true,
	}

	// Create user
	err := db.Create(&user).Error
	require.NoError(t, err)

	// Soft delete by setting IsActive to false
	err = db.Model(&user).Update("is_active", false).Error
	require.NoError(t, err)

	// Verify user is marked inactive but still exists
	var retrieved User
	err = db.First(&retrieved, "clerk_user_id = ?", user.ClerkUserID).Error
	require.NoError(t, err)

	assert.False(t, retrieved.IsActive)

	// User should still be queryable
	count := int64(0)
	err = db.Model(&User{}).Where("clerk_user_id = ?", user.ClerkUserID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestUserTableName(t *testing.T) {
	user := User{}
	tableName := user.TableName()
	assert.Equal(t, "users", tableName)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
