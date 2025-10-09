package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// VerifyAndLoadUser verifies a JWT token and loads the corresponding user from the database
// This is used for WebSocket authentication where the token is sent as a message payload
func VerifyAndLoadUser(ctx context.Context, db *gorm.DB, token string) (*store.User, error) {
	// Verify JWT with Clerk SDK
	sessionClaims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token:  token,
		Leeway: 10 * time.Second, // Allow 10 seconds clock skew
	})
	if err != nil {
		return nil, fmt.Errorf("jwt verification failed: %w", err)
	}

	// Extract Clerk user ID from claims
	clerkUserID := sessionClaims.Subject
	if clerkUserID == "" {
		return nil, fmt.Errorf("jwt claims missing 'sub' field")
	}

	// Load user from database
	var user store.User
	err = db.Where("clerk_user_id = ?", clerkUserID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found: %s", clerkUserID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive: %s", user.ID)
	}

	return &user, nil
}
