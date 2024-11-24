package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// Load and decode the encryption key
	// TODO Add this as a kube secret
	keyBase64, err := os.ReadFile("encryption.key")
	if err != nil {
		log.Fatalf("Failed to load encryption key: %v", err)
	}

	key, err := base64.StdEncoding.DecodeString(string(keyBase64))
	if err != nil {
		log.Fatalf("Failed to decode encryption key: %v", err)
	}

	if len(key) != 32 {
		log.Fatalf("Encryption key must be 32 bytes long, got %d bytes", len(key))
	}

	credentialService := &credential.CredentialService{
		DB:  db,
		Key: key,
	}

	// Set up the router
	r := mux.NewRouter()
	routes.CredentialRoutes(r, credentialService)

	// Set up a channel to listen for OS signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Start the API server in a goroutine
	server := &http.Server{Addr: ":8080", Handler: r}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting API server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	// Run the UI in the main goroutine
	go func() {
		<-ctx.Done()
		log.Println("UI closed. Stopping API server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to gracefully shut down API server: %v", err)
		}
	}()

	// Run the UI in the main thread
	log.Println("Starting UI")
	ui.RenderUI()

	// Wait for the API server goroutine to finish
	wg.Wait()
	log.Println("Application stopped gracefully")
}
