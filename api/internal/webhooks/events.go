package webhooks

import "encoding/json"

// WebhookEvent represents a Clerk webhook event
type WebhookEvent struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// ClerkUser represents user data from Clerk webhook
type ClerkUser struct {
	ID             string  `json:"id"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	ImageURL       *string `json:"image_url"`
	EmailAddresses []struct {
		EmailAddress string `json:"email_address"`
		ID           string `json:"id"`
	} `json:"email_addresses"`
}

// ExtractUser extracts user data from webhook payload
func ExtractUser(data map[string]interface{}) (*ClerkUser, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var user ClerkUser
	if err := json.Unmarshal(jsonData, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetPrimaryEmail returns the primary email address from the user
func (u *ClerkUser) GetPrimaryEmail() string {
	if len(u.EmailAddresses) > 0 {
		return u.EmailAddresses[0].EmailAddress
	}
	return ""
}
