package storage

import (
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
)

// StorageResource - интерфейс для хранилища освобождающий ресурсы
type StorageResource interface {
	// Resource
	Close() // освобождение ресурсов
}

// StorageOperators - интерфейс для хранилища операторов
type StorageOperators interface {
	// Users
	CreateOperator(user *models.Operator) (string, error) // сохранием нового пользователя
	GetAllOperators() ([]models.Operator, error)
	GetOperatorByUID(userUID string) (*models.Operator, error) // возвращает пользователя
	GetOperatorByName(username string) (*models.Operator, error)
}


// Storage - интерфейс для хранилища
type Storage interface {
	StorageResource
	StorageOperators
}