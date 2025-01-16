package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PerfectStepCoder/yp_go_nn/src/configs"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
)

type StorageInPostgres struct {
	connectionToDB string
	DB             *gorm.DB
}

func NewStorageInPostgres(connectionString string) (*StorageInPostgres, error) {

	newStorage := StorageInPostgres{connectionToDB: connectionString, DB: nil}

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Установка поддержки UID
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create extension: %v", err)
	}

	// Миграция схемы
	if err := db.AutoMigrate(&models.Operator{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	newStorage.DB = db

	return &newStorage, nil
}

func (s *StorageInPostgres) CreateOperator(user *models.Operator) (string, error) { // userUID
	if err := s.DB.Create(user).Error; err != nil {
		return "", err
	}
	return user.OperatorUID.String(), nil
}

func (s *StorageInPostgres) GetAllOperators() ([]models.Operator, error) {
	var users []models.Operator
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *StorageInPostgres) GetOperatorByUID(operatorUID string) (*models.Operator, error) {
	var user models.Operator
	if err := s.DB.Where("operator_uid = ?", operatorUID).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (s *StorageInPostgres) GetOperatorByName(name string) (*models.Operator, error) {
	var user models.Operator
	if err := s.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (s *StorageInPostgres) Close() {
	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err != nil {
			configs.Logger.Errorf("failed to get database connection: %s", err)
		}
		sqlDB.Close()
	}
}
