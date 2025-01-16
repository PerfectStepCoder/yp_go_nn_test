package security

import (
	"fmt"
	"time"

	"github.com/PerfectStepCoder/yp_go_nn/src/configs"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

type OperatorJWT struct {
	OperatorUID string
	Role    string
}

func GenerateJWT(operator *models.Operator) (string, error) {
	// Определение времени истечения токена
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создание токена с информацией о пользователе
	claims := &jwt.MapClaims{
		"operator_uid": operator.OperatorUID,
		"role":     operator.Role,
		"exp":      expirationTime.Unix(),
	}

	// Создаем токен с подписью
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.SettingsGlobal.SecretJWT)
}

// ExtractUserFromClaims извлекает данные пользователя из JWT claims
func ExtractUserFromClaims(claims *jwt.MapClaims) (*OperatorJWT, error) {

	operatorID, ok := (*claims)["operator_uid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_uid claim")
	}

	role, ok := (*claims)["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role claim")
	}

	// Возвращаем структуру пользователя
	return &OperatorJWT{
		OperatorUID: operatorID,
		Role:    role,
	}, nil
}
