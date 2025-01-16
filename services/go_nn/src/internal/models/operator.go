package models

import (
	"fmt"

	"github.com/google/uuid" // для поддержки массивов PostgreSQL
	"gorm.io/gorm"
)

type OperatorRole string

const (
	RoleUser  OperatorRole = "User"
	RoleAdmin OperatorRole = "Admin"
	RoleEngineer  OperatorRole = "Engineer"
)

type Operator struct {
	gorm.Model
	OperatorUID  uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();unique;not null;index" json:"operatorUID"`
	Name         string       `gorm:"unique;not null;index" json:"name"`
	Password     string       `gorm:"not null" json:"password"`
	Role         OperatorRole `gorm:"type:varchar(10);not null;default:'User';index" json:"role"`
}

func (o Operator) String() string {
	result := fmt.Sprintf("\tOperatorUID=%s\n\tName=%s\n\tPassword=%s\n\tRole=%s",
		o.OperatorUID.String(), o.Name, o.Password, o.Role,
	)
	return result
}

type OperatorRequest struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type OperatorResponseUID struct {
	OperatorUID string `json:"operator_uid"`
}

type OperatorResponse struct {
	OperatorUID  uuid.UUID `json:"operatorUID"`
	Name string    `json:"name"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

// Cтруктура для получения логина и пароля
type LoginRequest struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

// Cтруктура для ответа с JWT
type LoginResponse struct {
	Token string `json:"token"`
}
