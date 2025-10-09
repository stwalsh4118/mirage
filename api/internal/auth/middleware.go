package auth

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

const (
	contextKeyUser        = "auth:user"
	lastSeenThrottleHours = 1
)

var initOnce sync.Once

// InitClerk initializes the Clerk SDK with the secret key
func InitClerk(secretKey string) {
	initOnce.Do(func() {
		clerk.SetKey(secretKey)
		log.Info().Msg("clerk SDK initialized with secret key")
	})
}

// RequireAuth is middleware that verifies Clerk JWT and loads user
func RequireAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}
		log.Debug().Msgf("authHeader: %v", authHeader)

		// Parse "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		log.Debug().Msgf("parts: %v", parts)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format, expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}
		token := parts[1]

		// Verify JWT with Clerk SDK
		ctx := context.Background()
		sessionClaims, err := jwt.Verify(ctx, &jwt.VerifyParams{
			Token:  token,
			Leeway: 10 * time.Second, // Allow 10 seconds clock skew
		})
		if err != nil {
			log.Debug().Err(err).Msg("jwt verification failed")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract Clerk user ID from claims
		clerkUserID := sessionClaims.Subject
		if clerkUserID == "" {
			log.Error().Msg("jwt claims missing 'sub' field")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// Load user from database
		var user store.User
		err = db.Where("clerk_user_id = ?", clerkUserID).First(&user).Error
		if err == gorm.ErrRecordNotFound {
			log.Warn().
				Str("clerk_user_id", clerkUserID).
				Msg("user not found in database (webhook may not have processed yet)")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found - please try again in a moment",
			})
			c.Abort()
			return
		} else if err != nil {
			log.Error().Err(err).Str("clerk_user_id", clerkUserID).Msg("failed to load user")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to load user",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			log.Warn().Str("user_id", user.ID).Msg("inactive user attempted access")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "user account is inactive",
			})
			c.Abort()
			return
		}

		// Update LastSeenAt if needed (async, throttled)
		go updateLastSeenAt(db, &user)

		// Store user in context
		c.Set(contextKeyUser, &user)

		log.Debug().
			Str("user_id", user.ID).
			Str("email", user.Email).
			Msg("request authenticated")

		c.Next()
	}
}

// updateLastSeenAt updates user's LastSeenAt timestamp if last update was > 1 hour ago
func updateLastSeenAt(db *gorm.DB, user *store.User) {
	now := time.Now()

	// Throttle: only update if last seen was > 1 hour ago or never set
	if user.LastSeenAt != nil {
		elapsed := now.Sub(*user.LastSeenAt)
		if elapsed < lastSeenThrottleHours*time.Hour {
			return
		}
	}

	// Update in background (don't block request)
	err := db.Model(&store.User{}).
		Where("id = ?", user.ID).
		Update("last_seen_at", now).
		Error

	if err != nil {
		log.Warn().Err(err).Str("user_id", user.ID).Msg("failed to update last_seen_at")
	}
}
