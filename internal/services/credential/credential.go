package credential

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cooksey14/cred-manager/internal/models"
	sec "github.com/cooksey14/cred-manager/internal/security"
	"github.com/gorilla/mux"
)

type CredentialService struct {
	DB  *sql.DB
	Key []byte
}

// CreateCredential creates a new credential
func (s *CredentialService) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var cred models.Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	encryptedPassword, nonce, err := sec.Encrypt(cred.Password, s.Key)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO credentials (service, username, password, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)"
	result, err := s.DB.Exec(query, cred.Service, cred.Username, encryptedPassword+"|"+nonce)
	if err != nil {
		http.Error(w, "Failed to insert credential", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	cred.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cred)
}

// GetCredentials retrieves all credentials
func (s *CredentialService) GetCredentials(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("SELECT id, service, username, password, created_at FROM credentials")
	if err != nil {
		http.Error(w, "Failed to fetch credentials", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var creds []models.Credential
	for rows.Next() {
		var cred models.Credential
		var encryptedPassword string
		err := rows.Scan(&cred.ID, &cred.Service, &cred.Username, &encryptedPassword, &cred.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to read credentials", http.StatusInternalServerError)
			return
		}

		parts := strings.Split(encryptedPassword, "|")
		if len(parts) != 2 {
			http.Error(w, "Invalid password format", http.StatusInternalServerError)
			return
		}

		decryptedPassword, err := sec.Decrypt(parts[0], parts[1], s.Key)
		if err != nil {
			http.Error(w, "Failed to decrypt password", http.StatusInternalServerError)
			return
		}
		cred.Password = decryptedPassword

		creds = append(creds, cred)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(creds)
}

// UpdateCredential updates an existing credential
func (s *CredentialService) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var cred models.Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	encryptedPassword, nonce, err := sec.Encrypt(cred.Password, s.Key)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	query := "UPDATE credentials SET service = ?, username = ?, password = ? WHERE id = ?"
	_, err = s.DB.Exec(query, cred.Service, cred.Username, encryptedPassword+"|"+nonce, id)
	if err != nil {
		http.Error(w, "Failed to update credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCredential deletes a credential by ID
func (s *CredentialService) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := "DELETE FROM credentials WHERE id = ?"
	_, err := s.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCredential retrieves a single credential by ID
func (s *CredentialService) GetCredential(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var cred models.Credential

	query := "SELECT id, service, username, password, created_at FROM credentials WHERE id = ?"
	err := s.DB.QueryRow(query, id).Scan(&cred.ID, &cred.Service, &cred.Username, &cred.Password, &cred.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Credential not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve credential", http.StatusInternalServerError)
		return
	}

	parts := strings.Split(cred.Password, "|")
	if len(parts) != 2 {
		http.Error(w, "Invalid password format", http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := sec.Decrypt(parts[0], parts[1], s.Key)
	if err != nil {
		http.Error(w, "Failed to decrypt password", http.StatusInternalServerError)
		return
	}
	cred.Password = decryptedPassword

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cred)
}
