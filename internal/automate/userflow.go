package automate

import (
	"log"
	"net/http"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
)

func UserFlow() error {
	client := CreateClient()
	user := NewTestUser(true)
	count := 3

	if err := Register(client, user); err != nil {
		return err
	}

	if err := Login(client, user); err != nil {
		return err
	}

	if err := testFolderOperations(client, count); err != nil {
		return err
	}

	log.Println("All passed")
	return nil
}

func testFolderOperations(client *http.Client, count int) error {
	if _, err := GetFolders(client); err != nil {
		return err
	}

	if err := createMultipleFolders(client, count); err != nil {
		return err
	}

	folders, err := GetFolders(client)
	if err != nil {
		return err
	}

	if err := verifyAndUpdateFolders(client, folders); err != nil {
		return err
	}

	if err := testNoteOperations(client, folders, count); err != nil {
		return err
	}

	return deleteAllFolders(client, folders)
}

func createMultipleFolders(client *http.Client, count int) error {
	for i := 0; i < count; i++ {
		if err := CreateFolder(client); err != nil {
			return err
		}
	}
	return nil
}

func verifyAndUpdateFolders(client *http.Client, folders []models.FolderResponse) error {
	for _, folder := range folders {
		if err := GetFolder(client, folder.Slug); err != nil {
			return err
		}
	}

	for _, folder := range folders {
		if err := UpdateFolder(client, folder); err != nil {
			return err
		}
	}
	return nil
}

func testNoteOperations(client *http.Client, folders []models.FolderResponse, count int) error {
	if err := createNotesInFolders(client, folders, count); err != nil {
		return err
	}

	if _, err := GetNotes(client); err != nil {
		return err
	}

	if err := verifyNotesInFolders(client, folders); err != nil {
		return err
	}

	if err := updateNotesInFolders(client, folders); err != nil {
		return err
	}

	if _, err := GetNotes(client); err != nil {
		return err
	}

	return deleteNotesInFolders(client, folders)
}

func createNotesInFolders(client *http.Client, folders []models.FolderResponse, count int) error {
	for _, folder := range folders {
		for i := 0; i < count; i++ {
			if err := CreateNote(client, folder.Slug); err != nil {
				return err
			}
		}
	}
	return nil
}

func verifyNotesInFolders(client *http.Client, folders []models.FolderResponse) error {
	for _, folder := range folders {
		if _, err := GetNotesInFolder(client, folder.Slug); err != nil {
			return err
		}
	}

	for _, folder := range folders {
		notes, err := GetNotesInFolder(client, folder.Slug)
		if err != nil {
			return err
		}

		for _, note := range notes {
			if err := GetNoteInFolder(client, folder.Slug, note.Slug); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateNotesInFolders(client *http.Client, folders []models.FolderResponse) error {
	for _, folder := range folders {
		notes, err := GetNotesInFolder(client, folder.Slug)
		if err != nil {
			return err
		}

		for _, note := range notes {
			if err := UpdateNote(client, folder.Slug, note); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteNotesInFolders(client *http.Client, folders []models.FolderResponse) error {
	for _, folder := range folders {
		log.Println(">>> starting delete notes")
		notes, err := GetNotesInFolder(client, folder.Slug)
		if err != nil {
			return err
		}

		for _, note := range notes {
			if err := DeleteNote(client, folder.Slug, note.Slug); err != nil {
				return err
			}
		}
	}

	_, err := GetNotes(client)
	return err
}

func deleteAllFolders(client *http.Client, folders []models.FolderResponse) error {
	if _, err := GetFolders(client); err != nil {
		return err
	}

	for _, folder := range folders {
		if err := DeleteFolder(client, folder.Slug); err != nil {
			return err
		}
	}
	return nil
}
