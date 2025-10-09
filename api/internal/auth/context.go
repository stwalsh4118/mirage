package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/store"
)

var (
    ErrUserNotInContext = errors.New("user not found in request context")
)

// GetCurrentUser retrieves the authenticated user from the request context
func GetCurrentUser(c *gin.Context) (*store.User, error) {
    val, exists := c.Get(contextKeyUser)
    if !exists {
        return nil, ErrUserNotInContext
    }

    user, ok := val.(*store.User)
    if !ok {
        return nil, errors.New("invalid user type in context")
    }

    return user, nil
}

// GetCurrentUserID retrieves just the user ID from the request context
func GetCurrentUserID(c *gin.Context) (string, error) {
    user, err := GetCurrentUser(c)
    if err != nil {
        return "", err
    }
    return user.ID, nil
}

// MustGetCurrentUser retrieves the user or panics (use only when user is guaranteed)
func MustGetCurrentUser(c *gin.Context) *store.User {
    user, err := GetCurrentUser(c)
    if err != nil {
        panic("user not in context - did you forget to apply RequireAuth middleware?")
    }
    return user
}

