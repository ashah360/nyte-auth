package model

import "time"

type User struct {
	// General user information
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Role      string `json:"role" db:"role"`
	DiscordID string `json:"discord_id" db:"discord_id"`
	Upgraded  bool   `json:"upgraded" db:"upgraded"`

	// Timestamps
	AccessExpiresAt *time.Time `json:"access_expires_at" db:"access_expires_at"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
}

func NewUserWithDefaults() *User {
	return &User{
		Role: "USER",
	}
}
