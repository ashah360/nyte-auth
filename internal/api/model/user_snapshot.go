package model

import (
	"encoding/json"
	"time"
)

type UserSnapshot struct {
	ID              string     `json:"user_id"`           // User ID
	Token           string     `json:"token_id"`          // User token associated with the session
	Role            string     `json:"role"`              // User role
	Upgraded        bool       `json:"upgraded"`          // User upgraded status
	Grants          []string   `json:"grants"`            // User authorization grants
	AccessExpiresAt *time.Time `json:"access_expires_at"` // Expiration timestamp of the user's access to the bot
}

func (us *UserSnapshot) MarshalBinary() ([]byte, error) {
	return json.Marshal(us)
}

func (us *UserSnapshot) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, us)
}
