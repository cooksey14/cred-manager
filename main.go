package main

import (
	"log"
	"net/http"

	"github.com/cooksey14/cred-manager/internal/services/credential"
	"github.com/cooksey14/cred-manager/internal/services/credential/routes"
	"github.com/cooksey14/cred-manager/internal/store"
	"github.com/cooksey14/cred-manager/internal/ui"
	"github.com/gorilla/mux"
)

func main() {
	// Define the database file path and migrations directory
	dbFilePath := "./password_manager.db"
	migrationsDir := "./internal/store/migrations"

	// Initialize the database
	db := store.InitDatabase(dbFilePath, migrationsDir)
	defer db.Close()

	// Initialize the credential service
	credentialService := &credential.CredentialService{DB: db}

	// Set up the router
	r := mux.NewRouter()
	routes.RegisterCredentialRoutes(r, credentialService)

	// Start the UI
	ui.RenderUI(db)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
