package automate

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/brianvoe/gofakeit/v6"
)

func GetFolders(client *http.Client) error {
	url := UrlMap["get-folders"]

	_, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Folders fetched successfully")
	return nil
}

func CreateFolder(client *http.Client) error {
	url := UrlMap["create-folder"]

	newFolder := models.CreateFolderRequest{
		Name: fmt.Sprintf("%s-%d", gofakeit.Word(), gofakeit.Number(100, 999)),
	}

	_, err := DoJSONRequest(
		client,
		http.MethodPost,
		url,
		newFolder,
	)
	if err != nil {
		return err
	}

	log.Println("Folder created successfully")
	return nil
}
