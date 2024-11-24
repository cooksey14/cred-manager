package models

import (
	"database/sql"
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

type Vault struct {
	Entries []Credential `json:"entries"`
}
