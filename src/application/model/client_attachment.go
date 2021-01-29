package model

import (
	"gorm.io/gorm"
	"time"
)

type Attachment struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	StorageMethod string	`json:"storage_method"`
	Path string	`json:"path"`
	OwnerID uint	`json:"owner_id"`
	OwnerType string	`json:"owner_type"`
}