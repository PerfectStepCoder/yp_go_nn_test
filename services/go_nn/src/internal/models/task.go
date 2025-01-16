package models

import (
	"mime/multipart"

	"github.com/google/uuid" // для поддержки массивов PostgreSQL
)

type TaskOneRequest struct {
	TaskUID uuid.UUID `json:"taskUID"`
	Image   multipart.File `json:"image"` // Файл изображения
}

type TaskOneResponse struct {
	TaskUID string `json:"taskUID"`
	Label string
}
