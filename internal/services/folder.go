package services

import (
	"errors"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

var (
	ErrFolderNotFound = errors.New("folder not found")
	ErrFolderExists   = errors.New("folder already exists")
)

func GetFolders(db *gorm.DB, userId string) ([]models.FolderResponse, error) {
	var folders []models.FolderResponse

	if err := db.
		Model(&models.Folder{}).
		Select("id", "name", "slug").
		Where("user_id = ?", userId).
		Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

func GetFolder(db *gorm.DB, userId, folderSlug string) (models.FolderResponse, error) {
	var folder models.FolderResponse

	err := db.
		Model(&models.Folder{}).
		Select("id", "name", "slug").
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		First(&folder).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.FolderResponse{}, ErrFolderNotFound
		}
		return models.FolderResponse{}, err
	}

	return folder, nil
}

func CreateFolder(db *gorm.DB, body *models.CreateFolderRequest, userId string) error {
	var existing models.Folder

	err := db.
		Where("name = ? AND user_id = ?", body.Name, userId).
		First(&existing).Error

	if err == nil {
		return ErrFolderExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	folder := models.Folder{
		ID:     uuid.NewString(),
		UserID: userId,
		Name:   body.Name,
		Slug:   slug.Make(body.Name),
	}

	if err := db.Create(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrFolderExists
		}
		return err
	}

	return nil
}

func UpdateFolder(db *gorm.DB, body *models.UpdateFolderRequest, userId, folderSlug string) error {
	result := db.
		Model(&models.Folder{}).
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		Updates(map[string]any{
			"name": body.Name,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrFolderNotFound
	}

	return nil
}

func DeleteFolder(db *gorm.DB, userId, folderSlug string) error {
	result := db.
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		Delete(&models.Folder{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrFolderNotFound
	}

	return nil
}
