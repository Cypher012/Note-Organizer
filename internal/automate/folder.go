package automate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/brianvoe/gofakeit/v6"
)

func GetFolders(client *http.Client) (folders []models.FolderResponse, err error) {
	url := UrlMap["folders"]()

	body, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp := new(models.GetFoldersResponse)

	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	folders = resp.Folders

	log.Println("Folders fetched successfully")
	return folders, nil
}

func GetFolder(client *http.Client, folderSlug string) error {
	url := UrlMap["folders"](folderSlug)

	_, err := DoJSONRequest(
		client,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Folder fetched successfully")
	return nil
}

func CreateFolder(client *http.Client) error {
	url := UrlMap["folders"]()

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

func UpdateFolder(client *http.Client, folder models.FolderResponse) error {
	url := UrlMap["folders"](folder.Slug)

	UpdateFolder := models.UpdateFolderRequest{
		Name: strings.ToTitle(folder.Name),
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

	log.Println("Folder Updated successfully")
	return nil
}

func DeleteFolder(client *http.Client, folderSlug string) error {
	url := UrlMap["folders"](folderSlug)

	_, err := DoJSONRequest(
		client,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Folder Deleted successfully")
	return nil
}
