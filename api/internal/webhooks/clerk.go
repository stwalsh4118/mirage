package webhooks

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

const (
	EventUserCreated = "user.created"
	EventUserUpdated = "user.updated"
	EventUserDeleted = "user.deleted"
)

const defaultWebhookTTL = 24 * time.Hour

// ClerkWebhookHandler handles Clerk webhook events
type ClerkWebhookHandler struct {
	db            *gorm.DB
	webhookSecret string
	idCache       *WebhookIDCache
}

// NewClerkWebhookHandler creates a new Clerk webhook handler
func NewClerkWebhookHandler(db *gorm.DB, webhookSecret string) *ClerkWebhookHandler {
	return &ClerkWebhookHandler{
		db:            db,
		webhookSecret: webhookSecret,
		idCache:       NewWebhookIDCache(defaultWebhookTTL),
	}
}

// HandleWebhook processes incoming Clerk webhook events
func (h *ClerkWebhookHandler) HandleWebhook(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read webhook body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Verify webhook signature using Svix
	if err := VerifyWebhookSignature(body, c.Request.Header, h.webhookSecret); err != nil {
		log.Warn().Err(err).Msg("webhook signature verification failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// Parse webhook event
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().Err(err).Msg("failed to parse webhook payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook payload"})
		return
	}

	// Get webhook ID from Svix header (this is the unique webhook delivery ID)
	webhookID := c.GetHeader("svix-id")
	if webhookID == "" {
		log.Warn().Msg("webhook missing svix-id header")
		// Don't fail - just process without idempotency check
		webhookID = event.Type + "-" + time.Now().Format(time.RFC3339Nano)
	}

	// Check for duplicate webhooks (idempotency)
	if h.idCache.IsProcessed(webhookID) {
		log.Info().Str("webhook_id", webhookID).Msg("webhook already processed, skipping")
		c.JSON(http.StatusOK, gin.H{"status": "already processed"})
		return
	}

	// Process event based on type
	switch event.Type {
	case EventUserCreated:
		err = h.handleUserCreated(event)
	case EventUserUpdated:
		err = h.handleUserUpdated(event)
	case EventUserDeleted:
		err = h.handleUserDeleted(event)
	default:
		log.Warn().Str("event_type", event.Type).Msg("unsupported webhook event type")
		c.JSON(http.StatusOK, gin.H{"status": "unsupported event type"})
		return
	}

	if err != nil {
		log.Error().Err(err).Str("event_type", event.Type).Msg("failed to process webhook event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process webhook"})
		return
	}

	// Mark webhook as processed
	h.idCache.MarkProcessed(webhookID)

	log.Info().
		Str("webhook_id", webhookID).
		Str("event_type", event.Type).
		Msg("webhook processed successfully")

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handleUserCreated creates a new user in the database
func (h *ClerkWebhookHandler) handleUserCreated(event WebhookEvent) error {
	clerkUser, err := ExtractUser(event.Data)
	if err != nil {
		return err
	}

	email := clerkUser.GetPrimaryEmail()
	if email == "" {
		log.Warn().Str("clerk_user_id", clerkUser.ID).Msg("user has no email address")
		// Don't fail the webhook, just log the warning
		return nil
	}

	// Check if user already exists (idempotency at DB level)
	var existingUser store.User
	err = h.db.Where("clerk_user_id = ?", clerkUser.ID).First(&existingUser).Error
	if err == nil {
		log.Info().Str("clerk_user_id", clerkUser.ID).Msg("user already exists, skipping creation")
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create new user
	user := store.User{
		ID:              uuid.New().String(),
		ClerkUserID:     clerkUser.ID,
		Email:           email,
		FirstName:       clerkUser.FirstName,
		LastName:        clerkUser.LastName,
		ProfileImageURL: clerkUser.ImageURL,
		IsActive:        true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		return err
	}

	log.Info().
		Str("user_id", user.ID).
		Str("clerk_user_id", user.ClerkUserID).
		Str("email", user.Email).
		Msg("user created from webhook")

	return nil
}

// handleUserUpdated updates an existing user in the database
func (h *ClerkWebhookHandler) handleUserUpdated(event WebhookEvent) error {
	clerkUser, err := ExtractUser(event.Data)
	if err != nil {
		return err
	}

	email := clerkUser.GetPrimaryEmail()
	if email == "" {
		log.Warn().Str("clerk_user_id", clerkUser.ID).Msg("user has no email address")
		return nil
	}

	// Find existing user
	var user store.User
	err = h.db.Where("clerk_user_id = ?", clerkUser.ID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		log.Warn().Str("clerk_user_id", clerkUser.ID).Msg("user not found for update, creating instead")
		// Create the user if it doesn't exist (race condition or missed webhook)
		return h.handleUserCreated(event)
	} else if err != nil {
		return err
	}

	// Update user fields
	updates := map[string]interface{}{
		"email":             email,
		"first_name":        clerkUser.FirstName,
		"last_name":         clerkUser.LastName,
		"profile_image_url": clerkUser.ImageURL,
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		return err
	}

	log.Info().
		Str("user_id", user.ID).
		Str("clerk_user_id", user.ClerkUserID).
		Msg("user updated from webhook")

	return nil
}

// handleUserDeleted hard-deletes a user from the database
func (h *ClerkWebhookHandler) handleUserDeleted(event WebhookEvent) error {
	clerkUser, err := ExtractUser(event.Data)
	if err != nil {
		return err
	}

	// Find existing user
	var user store.User
	err = h.db.Where("clerk_user_id = ?", clerkUser.ID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		log.Warn().Str("clerk_user_id", clerkUser.ID).Msg("user not found for deletion")
		// Not an error - user might have been deleted already
		return nil
	} else if err != nil {
		return err
	}

	// Hard delete the user record
	if err := h.db.Delete(&user).Error; err != nil {
		return err
	}

	log.Info().
		Str("user_id", user.ID).
		Str("clerk_user_id", user.ClerkUserID).
		Msg("user deleted from webhook")

	return nil
}
