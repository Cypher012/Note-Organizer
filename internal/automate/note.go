package automate

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/brianvoe/gofakeit/v6"
)

type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	FolderID string `json:"folder_id"`
}

func GetNotes(client *http.Client) (notes []models.NoteResponse, err error) {
	url := UrlMap["get-notes"]()

	body, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp := new(models.GetNotesResponse)

	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	notes = resp.Notes

	log.Println("Notes fetched successfully")
	return notes, nil
}

func GetNotesInFolder(client *http.Client, folderSlug string) (notes []models.NoteResponse, err error) {
	url := UrlMap["folders"](folderSlug, "notes")
	body, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp := new(models.GetNotesResponse)

	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	notes = resp.Notes

	log.Println("Notes in folder fetched successfully: ", folderSlug)
	return notes, nil
}

func GetNoteInFolder(client *http.Client, folderSlug string, noteSlug string) error {
	url := UrlMap["folders"](folderSlug, "notes", noteSlug)
	_, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Note in folder fetched successfully: ", folderSlug, noteSlug)
	return nil
}

func CreateNote(client *http.Client, folderSlug string) error {
	url := UrlMap["folders"](folderSlug, "notes")

	var content = gofakeit.Paragraph(1, 2, 10, " ")

	note := models.CreateNoteRequest{
		Name:    gofakeit.Sentence(3),
		Content: &content,
	}

	_, err := DoJSONRequest(
		client,
		http.MethodPost,
		url,
		note,
	)
	if err != nil {
		return err
	}

	log.Println("Note created for folder:", folderSlug)
	return nil
}

func UpdateNote(client *http.Client, folderSlug string, note models.NoteResponse) error {
	url := UrlMap["folders"](folderSlug, "notes", note.Slug)

	UpdateFolder := models.UpdateNoteRequest{
		Name:    strings.ToTitle(note.Name),
		Content: strPtr(strings.ToTitle(*note.Content)),
	}

	_, err := DoJSONRequest(
		client,
		http.MethodPut,
		url,
		UpdateFolder,
	)
	if err != nil {
		return err
	}

	log.Println("Note Updated successfully: ", folderSlug, note.Slug)
	return nil
}

func DeleteNote(client *http.Client, folderSlug string, noteSlug string) error {
	url := UrlMap["folders"](folderSlug, "notes", noteSlug)

	_, err := DoJSONRequest(
		client,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Note Deleted successfully: ", noteSlug)
	return nil
}
