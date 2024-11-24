package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/cooksey14/cred-manager/internal/models"
)

// API Base URL
const apiBaseURL = "http://localhost:8080"

// Fetch credentials via API
func fetchCredentials() ([]models.Credential, error) {
	resp, err := http.Get(fmt.Sprintf("%s/credentials", apiBaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch credentials: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch credentials: status code %d", resp.StatusCode)
	}

	var creds []models.Credential
	if err := json.NewDecoder(resp.Body).Decode(&creds); err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %v", err)
	}

	return creds, nil
}

// Delete a credential via API
func deleteCredential(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/credentials/%d", apiBaseURL, id), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete credential: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete credential: %s", body)
	}

	return nil
}

// Add a new credential via API
func addCredential(cred models.Credential) error {
	data, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/credentials", apiBaseURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to add credential: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add credential: %s", body)
	}

	return nil
}

// RenderUI creates and displays the main UI window
func RenderUI() {
	a := app.New()
	w := a.NewWindow("Password Manager")

	// Track the selected index
	var selectedIndex int = -1

	// Fetch and display credentials
	credentialsList := widget.NewList(
		func() int {
			creds, _ := fetchCredentials()
			return len(creds)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			creds, _ := fetchCredentials()
			item.(*widget.Label).SetText(fmt.Sprintf("%s (%s)", creds[id].Service, creds[id].Username))
		},
	)

	credentialsList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
	}

	// Add a new credential
	showAddCredentialDialog := func(w fyne.Window, credentialsList *widget.List) {
		serviceEntry := widget.NewEntry()
		serviceEntry.SetPlaceHolder("Service Name")
		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("Username")
		passwordEntry := widget.NewEntry()
		passwordEntry.SetPlaceHolder("Password")
		passwordEntry.Password = true

		dialog.ShowCustomConfirm("Add Credential", "Save", "Cancel",
			container.NewVBox(
				widget.NewLabel("Service:"),
				serviceEntry,
				widget.NewLabel("Username:"),
				usernameEntry,
				widget.NewLabel("Password:"),
				passwordEntry,
			),
			func(confirm bool) {
				if confirm {
					service := serviceEntry.Text
					username := usernameEntry.Text
					password := passwordEntry.Text

					if service == "" || username == "" || password == "" {
						dialog.ShowError(fmt.Errorf("all fields are required"), w)
						return
					}

					cred := models.Credential{
						Service:  service,
						Username: username,
						Password: password,
					}

					if err := addCredential(cred); err != nil {
						dialog.ShowError(fmt.Errorf("failed to add credential: %v", err), w)
						return
					}

					credentialsList.Refresh()
				}
			},
			w,
		)
	}

	// Delete a selected credential
	deleteButton := widget.NewButton("Delete Selected", func() {
		if selectedIndex < 0 {
			dialog.ShowError(fmt.Errorf("no credential selected"), w)
			return
		}

		creds, _ := fetchCredentials()
		if selectedIndex >= len(creds) {
			dialog.ShowError(fmt.Errorf("invalid selection"), w)
			return
		}

		cred := creds[selectedIndex]
		if err := deleteCredential(cred.ID); err != nil {
			dialog.ShowError(fmt.Errorf("failed to delete credential: %v", err), w)
			return
		}

		credentialsList.Unselect(selectedIndex)
		credentialsList.Refresh()
	})

	// Refresh credentials list
	refreshButton := widget.NewButton("Refresh", func() {
		credentialsList.Refresh()
	})

	// Add a credential
	addButton := widget.NewButton("Add Credential", func() {
		showAddCredentialDialog(w, credentialsList)
	})

	// Layout the UI
	w.SetContent(container.NewBorder(
		container.NewVBox(refreshButton, addButton, deleteButton),
		nil, nil, nil,
		credentialsList,
	))

	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}
