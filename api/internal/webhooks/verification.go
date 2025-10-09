package webhooks

import (
	"errors"
	"net/http"

	svix "github.com/svix/svix-webhooks/go"
)

var (
	ErrInvalidSignature = errors.New("invalid webhook signature")
	ErrMissingSignature = errors.New("missing webhook signature header")
)

// VerifyWebhookSignature verifies the Clerk webhook signature using Svix SDK
func VerifyWebhookSignature(payload []byte, headers http.Header, secret string) error {
	wh, err := svix.NewWebhook(secret)
	if err != nil {
		return err
	}

	// Verify the webhook signature
	// Svix will check svix-id, svix-timestamp, and svix-signature headers
	err = wh.Verify(payload, headers)
	if err != nil {
		return ErrInvalidSignature
	}

	return nil
}
