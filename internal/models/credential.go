package models

import (
	"database/sql"
	"time"
)

type Credential struct {
	ID        int    `json:"id"`
	Service   string `json:"service"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// CredentialService handles operations on credentials
type CredentialService struct {
	DB *sql.DB
}

type PasswordEntry struct {
	Service   string    `json:"service"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Vault struct {
	Entries []PasswordEntry `json:"entries"`
}
