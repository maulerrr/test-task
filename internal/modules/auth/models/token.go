package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid" json:"user_id"`
	RefreshTokenHash string    `gorm:"size:255" json:"-"`
	IPAddress        string    `gorm:"size:45" json:"ip_address"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}
