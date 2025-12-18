package models

import "time"

type GenderType string
type RoleType string

const (
	Male        GenderType = "male"
	Female      GenderType = "female"
	Unspecified GenderType = "unspecified"
)

const (
	RoleUser       RoleType = "user"
	RoleAdmin      RoleType = "admin"
	RoleSuperAdmin RoleType = "super_admin"
)

type User struct {
	ID        string     `json:"id" gorm:"primaryKey;not null"`
	Username  string     `json:"username" gorm:"not null"`
	Email     string     `json:"email" gorm:"not null;uniqueIndex"`
	Password  string     `json:"-" gorm:"not null"`
	Gender    GenderType `json:"gender" gorm:"default:'unspecified'"`
	Role      RoleType   `json:"role" gorm:"default:'user'"`
	Bio       string     `json:"bio"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	Username string     `json:"username" validate:"omitempty,min=3"`
	Gender   GenderType `json:"gender" validate:"omitempty,oneof=male female unspecified"`
	Bio      string     `json:"bio" validate:"omitempty,min=10"`
}
