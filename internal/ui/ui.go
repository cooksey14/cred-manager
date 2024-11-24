package ui

import (
	"database/sql"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/cooksey14/cred-manager/internal/models"
)

// RenderUI creates and displays the main UI window
func RenderUI(db *sql.DB) {
	// Create a new Fyne application
	a := app.New()

	// Create a new window
	w := a.NewWindow("Password Manager")

	// Fetch credentials from the database
	fetchCredentials := func() ([]models.Credential, error) {
		rows, err := db.Query("SELECT id, service, username, password, created_at FROM credentials")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var creds []models.Credential
		for rows.Next() {
			var cred models.Credential
			err := rows.Scan(&cred.ID, &cred.Service, &cred.Username, &cred.Password, &cred.CreatedAt)
			if err != nil {
				return nil, err
			}
			creds = append(creds, cred)
		}
		return creds, nil
	}

	// Delete a credential from the database
	deleteCredential := func(id int) error {
		query := "DELETE FROM credentials WHERE id = ?"
		_, err := db.Exec(query, id)
		return err
	}

	// Show a dialog to add a new credential
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

					// Insert into the database
					query := "INSERT INTO credentials (service, username, password) VALUES (?, ?, ?)"
					_, err := db.Exec(query, service, username, password)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to add credential: %v", err), w)
						return
					}

					credentialsList.Refresh()
				}
			},
			w,
		)
	}

	// Track the selected index
	var selectedIndex int = -1

	// Create the list widget
	credentialsList := widget.NewList(
		func() int {
			creds, _ := fetchCredentials()
			return len(creds)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Template for each row
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			creds, _ := fetchCredentials()
			item.(*widget.Label).SetText(fmt.Sprintf("%s (%s)", creds[id].Service, creds[id].Username))
		},
	)

	// Set the OnSelected callback
	credentialsList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
	}

	// Add button to refresh the list
	refreshButton := widget.NewButton("Refresh", func() {
		credentialsList.Refresh()
	})

	// Add button to add a new credential
	addButton := widget.NewButton("Add Credential", func() {
		showAddCredentialDialog(w, credentialsList)
	})

	// Add button to delete a selected credential
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

	// Arrange buttons and list in the UI
	w.SetContent(container.NewBorder(
		container.NewVBox(refreshButton, addButton, deleteButton),
		nil, nil, nil,
		credentialsList,
	))

	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}

// func copyToClipboard(app fyne.App, content string) {
// 	var clipboard string
// 	clipboard = app.Clipboard()
// 	clipboard.SetContent(content)

// 	go func() {
// 		time.Sleep(30 * time.Second)
// 		clipboard.SetContent("")
// 	}()
// }
