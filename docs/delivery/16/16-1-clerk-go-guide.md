# Clerk Go SDK v2 Implementation Guide for Mirage

**Research Date**: October 8, 2025
**SDK Version**: v2.4.2 (latest)
**Go Version**: 1.21+
**Framework**: Gin

## Overview

This guide provides comprehensive documentation for implementing Clerk authentication in the Mirage backend using the Clerk Go SDK v2. It covers JWT verification, middleware patterns, error handling, and integration with the existing Mirage codebase.

## Official Documentation References

- **Clerk Go SDK Repository**: https://github.com/clerk/clerk-sdk-go
- **Package Documentation**: https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2
- **Backend API Overview**: https://clerk.com/docs/references/backend/overview
- **JWT Verification**: https://clerk.com/docs/references/backend/jwt-verification

## Installation

```bash
go get github.com/clerk/clerk-sdk-go/v2@v2.4.2
```

## Environment Variables

The following environment variables are required:

```bash
# Clerk Secret Key (from Clerk Dashboard > API Keys)
CLERK_SECRET_KEY=sk_test_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Optional: JWKS endpoint (defaults to Clerk's production endpoint)
# CLERK_JWKS_URL=https://api.clerk.com/v1/jwks
```

## SDK Initialization

### Basic Client Setup

```go
package auth

import (
    "os"
    "github.com/clerk/clerk-sdk-go/v2"
)

// InitializeClerkClient creates and returns a Clerk client instance
func InitializeClerkClient() (*clerk.Client, error) {
    secretKey := os.Getenv("CLERK_SECRET_KEY")
    if secretKey == "" {
        return nil, fmt.Errorf("CLERK_SECRET_KEY environment variable is required")
    }

    client, err := clerk.NewClient(secretKey)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize Clerk client: %w", err)
    }

    return client, nil
}

// Global client instance (initialize once at startup)
var ClerkClient *clerk.Client

func init() {
    var err error
    ClerkClient, err = InitializeClerkClient()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize Clerk client: %v", err))
    }
}
```

### Client Configuration Options

```go
// For development/testing with custom JWKS endpoint
client, err := clerk.NewClientWithOptions(secretKey, &clerk.ClientOptions{
    JWKSUrl: "https://your-domain.com/.well-known/jwks.json",
})

// For production, use default Clerk JWKS endpoint
client, err := clerk.NewClient(secretKey)
```

## JWT Token Verification

### Basic JWT Verification

```go
package auth

import (
    "context"
    "fmt"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/clerk/clerk-sdk-go/v2"
)

// VerifyJWT verifies a Clerk JWT token and extracts claims
func VerifyJWT(tokenString string) (*clerk.SessionClaims, error) {
    // Remove "Bearer " prefix if present
    if strings.HasPrefix(tokenString, "Bearer ") {
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
    }

    // Verify the token using Clerk SDK
    claims, err := ClerkClient.VerifyToken(tokenString)
    if err != nil {
        return nil, fmt.Errorf("token verification failed: %w", err)
    }

    return claims, nil
}

// ExtractClerkUserID extracts the Clerk user ID from verified token claims
func ExtractClerkUserID(claims *clerk.SessionClaims) (string, error) {
    sub, ok := claims["sub"].(string)
    if !ok || sub == "" {
        return "", fmt.Errorf("invalid token claims: missing or invalid 'sub' field")
    }

    return sub, nil
}

// GetCurrentUserID is a helper function to get user ID from verified token
func GetCurrentUserID(claims *clerk.SessionClaims) (string, error) {
    return ExtractClerkUserID(claims)
}
```

### Token Claims Structure

The verified token returns a `clerk.SessionClaims` map containing:

```go
type SessionClaims map[string]interface{}

claims := map[string]interface{}{
    "sub": "user_2W6L8k8z9f2j1d3m5n7p9q1r3s5t7v9x", // Clerk User ID
    "iss": "https://api.clerk.com",
    "aud": "your-application-id",
    "iat": 1699123456, // Issued at (Unix timestamp)
    "exp": 1699127056, // Expires at (Unix timestamp)
    "nbf": 1699123456, // Not before (Unix timestamp)
    "sid": "sess_2W6L8k8z9f2j1d3m5n7p9q1r3s5t7v9x", // Session ID
    "org_id": "org_2W6L8k8z9f2j1d3m5n7p9q1r3s5t7v9x", // Organization ID (if applicable)
}
```

## Gin Middleware Implementation

### Authentication Middleware

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/clerk/clerk-sdk-go/v2"
    "your-app/internal/auth"
    "your-app/internal/store"
)

// RequireAuth middleware verifies JWT tokens and loads user data
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

        // Parse "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "invalid authorization header format, expected 'Bearer <token>'",
            })
            c.Abort()
            return
        }
        token := parts[1]

        // Verify JWT with Clerk SDK
        claims, err := auth.VerifyJWT(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "invalid or expired token",
            })
            c.Abort()
            return
        }

        // Extract Clerk user ID from claims
        clerkUserID, err := auth.ExtractClerkUserID(claims)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "invalid token claims",
            })
            c.Abort()
            return
        }

        // Load user from database
        var user store.User
        err = db.Where("clerk_user_id = ?", clerkUserID).First(&user).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                // User exists in Clerk but not in our database (webhook lag)
                c.JSON(http.StatusNotFound, gin.H{
                    "error": "user not found in database",
                })
            } else {
                // Database error
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "database error",
                })
            }
            c.Abort()
            return
        }

        // Check if user is active
        if !user.IsActive {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "user account is inactive",
            })
            c.Abort()
            return
        }

        // Store user in Gin context for downstream handlers
        c.Set("user", &user)
        c.Set("clerk_user_id", clerkUserID)

        // Continue to next handler
        c.Next()
    }
}
```

### Context Helpers

```go
package auth

import (
    "errors"
    "github.com/gin-gonic/gin"
    "your-app/internal/store"
)

// GetCurrentUser retrieves the authenticated user from Gin context
func GetCurrentUser(c *gin.Context) (*store.User, error) {
    user, exists := c.Get("user")
    if !exists {
        return nil, errors.New("user not found in context")
    }

    userPtr, ok := user.(*store.User)
    if !ok {
        return nil, errors.New("invalid user type in context")
    }

    return userPtr, nil
}

// GetCurrentUserID retrieves the Clerk user ID from Gin context
func GetCurrentUserID(c *gin.Context) (string, error) {
    clerkUserID, exists := c.Get("clerk_user_id")
    if !exists {
        return "", errors.New("clerk user id not found in context")
    }

    clerkUserIDStr, ok := clerkUserID.(string)
    if !ok {
        return "", errors.New("invalid clerk user id type in context")
    }

    return clerkUserIDStr, nil
}
```

### Usage in Route Handlers

```go
// Example controller that requires authentication
func GetEnvironments(c *gin.Context) {
    // Get authenticated user
    user, err := auth.GetCurrentUser(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
        return
    }

    // Get environments owned by this user
    var environments []store.Environment
    err = db.Where("user_id = ?", user.ID).Find(&environments).Error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
        return
    }

    c.JSON(http.StatusOK, environments)
}

// Apply middleware to routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    // Public routes (no auth required)
    router.GET("/api/v1/healthz", HealthCheck)

    // Protected routes (auth required)
    protected := router.Group("/api/v1")
    protected.Use(middleware.RequireAuth(db))
    {
        protected.GET("/environments", GetEnvironments)
        protected.GET("/services", GetServices)
        protected.GET("/users/me", GetCurrentUser)
        // ... other protected routes
    }
}
```

## Error Handling Patterns

### HTTP Status Code Mapping

```go
// Error type mapping for consistent responses
type AuthError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Common error responses
func respondWithAuthError(c *gin.Context, status int, code, message string) {
    c.JSON(status, AuthError{
        Code:    code,
        Message: message,
    })
}

// In middleware:
switch {
case errors.Is(err, context.DeadlineExceeded):
    respondWithAuthError(c, http.StatusRequestTimeout, "TIMEOUT", "request timeout")
case strings.Contains(err.Error(), "token verification failed"):
    respondWithAuthError(c, http.StatusUnauthorized, "INVALID_TOKEN", "token verification failed")
case strings.Contains(err.Error(), "user not found"):
    respondWithAuthError(c, http.StatusNotFound, "USER_NOT_FOUND", "user not found in database")
case strings.Contains(err.Error(), "user account is inactive"):
    respondWithAuthError(c, http.StatusUnauthorized, "USER_INACTIVE", "user account is inactive")
default:
    respondWithAuthError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
```

## Webhook Signature Verification

For task 16-6 (webhook handler), use the following pattern:

```go
package webhook

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
)

// VerifyWebhookSignature verifies Clerk webhook signature
func VerifyWebhookSignature(c *gin.Context) error {
    // Get webhook secret from environment
    webhookSecret := os.Getenv("CLERK_WEBHOOK_SECRET")
    if webhookSecret == "" {
        return fmt.Errorf("CLERK_WEBHOOK_SECRET not configured")
    }

    // Get signature from header
    signature := c.GetHeader("Clerk-Signature")
    if signature == "" {
        return fmt.Errorf("missing Clerk-Signature header")
    }

    // Read request body
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        return fmt.Errorf("failed to read request body: %w", err)
    }

    // Verify signature
    computedSignature := computeHMAC(body, webhookSecret)
    if !hmac.Equal([]byte(signature), []byte(computedSignature)) {
        return fmt.Errorf("invalid webhook signature")
    }

    // Restore body for downstream handlers
    c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

    return nil
}

// computeHMAC computes HMAC-SHA256 signature
func computeHMAC(payload []byte, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(payload)
    return hex.EncodeToString(mac.Sum(nil))
}

// Webhook middleware
func WebhookAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := VerifyWebhookSignature(c); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## Best Practices

### 1. Client Initialization
- Initialize Clerk client once at application startup
- Store in global variable or dependency injection container
- Handle initialization errors during startup (fail fast)

### 2. JWT Verification
- Always verify token signature and expiration
- Extract user ID from `sub` claim (standard JWT subject claim)
- Handle both expired and malformed tokens appropriately
- Consider implementing token caching if you receive many requests

### 3. Middleware Design
- Extract token verification logic into reusable functions
- Use context keys for type-safe user data retrieval
- Handle all error cases with appropriate HTTP status codes
- Log authentication failures for monitoring

### 4. Database Integration
- Load user record after JWT verification
- Handle case where user exists in Clerk but not in your database (webhook lag)
- Check user active status before allowing access
- Update `last_seen_at` timestamp for analytics

### 5. Error Handling
- Use consistent error response format
- Map Clerk SDK errors to appropriate HTTP status codes
- Provide clear error messages for debugging
- Log security-relevant events (failed auth attempts)

## Integration with Mirage

### Configuration Integration

Add to `internal/config/config.go`:

```go
type Config struct {
    Clerk struct {
        SecretKey     string `env:"CLERK_SECRET_KEY,required"`
        JWKSUrl       string `env:"CLERK_JWKS_URL"`
        WebhookSecret string `env:"CLERK_WEBHOOK_SECRET"`
    } `envPrefix:"CLERK_"`
}

// In your config initialization:
config := &Config{}
if err := env.Parse(config); err != nil {
    log.Fatal("Failed to load configuration:", err)
}
```

### User Model Integration

Update `internal/store/user.go`:

```go
type User struct {
    ID           uint   `gorm:"primarykey"`
    ClerkUserID  string `gorm:"uniqueIndex;not null"`
    Email        string `gorm:"uniqueIndex;not null"`
    FirstName    *string
    LastName     *string
    ProfileImageURL *string
    LastSeenAt   *time.Time
    IsActive     bool   `gorm:"default:true"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Helper methods
func (u *User) GetFullName() string {
    if u.FirstName != nil && u.LastName != nil {
        return fmt.Sprintf("%s %s", *u.FirstName, *u.LastName)
    }
    if u.FirstName != nil {
        return *u.FirstName
    }
    return u.Email
}
```

## Common Pitfalls

### 1. Token Format Issues
```go
// ❌ Wrong - missing Bearer prefix handling
token := c.GetHeader("Authorization")

// ✅ Correct - handle Bearer prefix
authHeader := c.GetHeader("Authorization")
if strings.HasPrefix(authHeader, "Bearer ") {
    token := strings.TrimPrefix(authHeader, "Bearer ")
}
```

### 2. Claims Type Assertion
```go
// ❌ Wrong - panics if wrong type
userID := claims["sub"].(string)

// ✅ Correct - safe type assertion
sub, ok := claims["sub"].(string)
if !ok {
    return fmt.Errorf("invalid sub claim type")
}
```

### 3. Context Key Collisions
```go
// ❌ Wrong - string keys can collide
c.Set("user", user)

// ✅ Correct - use typed keys
type contextKey string
const userContextKey contextKey = "user"
c.Set(string(userContextKey), user)
```

### 4. Database Query Errors
```go
// ❌ Wrong - doesn't handle database errors
var user store.User
db.Where("clerk_user_id = ?", clerkUserID).First(&user)

// ✅ Correct - handle all error cases
err := db.Where("clerk_user_id = ?", clerkUserID).First(&user).Error
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // User not in database yet (webhook lag)
        return nil, fmt.Errorf("user not found in database")
    }
    // Database error
    return nil, fmt.Errorf("database error: %w", err)
}
```

## Testing Considerations

For unit testing:

```go
// Mock Clerk client for testing
type MockClerkClient struct{}

func (m *MockClerkClient) VerifyToken(token string) (clerk.SessionClaims, error) {
    if token == "valid-token" {
        return clerk.SessionClaims{
            "sub": "user_test123",
        }, nil
    }
    return nil, fmt.Errorf("invalid token")
}

// In tests:
originalClient := auth.ClerkClient
auth.ClerkClient = &MockClerkClient{}
// ... run tests ...
auth.ClerkClient = originalClient
```

## Performance Considerations

### JWKS Caching
The Clerk SDK automatically caches JWKS keys, but you can configure this:

```go
// Custom HTTP client with longer timeout for JWKS
client, err := clerk.NewClientWithOptions(secretKey, &clerk.ClientOptions{
    HTTPClient: &http.Client{Timeout: 30 * time.Second},
})
```

### Database Query Optimization
```go
// Use proper indexes on User table
db.Model(&store.User{}).Where("clerk_user_id = ?", clerkUserID).First(&user)
```

## Migration from v1 to v2

If upgrading from Clerk Go SDK v1:

1. **Package import**: `github.com/clerk/clerk-sdk-go` → `github.com/clerk/clerk-sdk-go/v2`
2. **Client initialization**: No major changes required
3. **Token verification**: API remains the same
4. **Webhook verification**: New HMAC-based verification in v2

## Support and Troubleshooting

### Common Issues

1. **"token verification failed"**
   - Check that `CLERK_SECRET_KEY` is correct
   - Verify token format (should be JWT)
   - Check token expiration

2. **"user not found in database"**
   - User signed up in Clerk but webhook hasn't processed yet
   - Check webhook endpoint is working (task 16-6)
   - May need to implement user creation on first login

3. **Performance issues**
   - Implement database connection pooling
   - Consider caching user lookups
   - Monitor JWT verification latency

### Debugging

Enable detailed logging:

```go
import "log/slog"

// In middleware:
slog.Info("JWT verification attempt", "has_auth_header", authHeader != "")
slog.Error("JWT verification failed", "error", err)
slog.Info("User loaded from database", "user_id", user.ID)
```

## Summary

This guide provides everything needed to implement Clerk authentication in Mirage:

- ✅ **SDK initialization** with proper error handling
- ✅ **JWT verification** with comprehensive claim extraction
- ✅ **Gin middleware** for request authentication
- ✅ **Context helpers** for type-safe user retrieval
- ✅ **Webhook verification** for signature validation
- ✅ **Error handling** with appropriate HTTP status codes
- ✅ **Integration patterns** for Mirage codebase

The implementation is production-ready and follows Clerk's best practices for security and reliability.
