package routes

import (
	"github.com/cooksey14/cred-manager/internal/services/credential"
	"github.com/gorilla/mux"
)

// CredentialRoutes sets up routes for the credential service
func CredentialRoutes(r *mux.Router, service *credential.CredentialService) {
	r.HandleFunc("/credentials", service.CreateCredential).Methods("POST")
	r.HandleFunc("/credentials", service.GetCredentials).Methods("GET")
	r.HandleFunc("/credentials/{id}", service.GetCredential).Methods("GET")
	r.HandleFunc("/credentials/{id}", service.UpdateCredential).Methods("PUT")
	r.HandleFunc("/credentials/{id}", service.DeleteCredential).Methods("DELETE")
}
