package httpp // http protocol

import (
	"context"
	"fmt"
	"github.com/PerfectStepCoder/yp_go_nn/src/configs"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/security"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

// ContextKey тип для ключей в контексте, чтобы избежать конфликтов
type ContextKey string

const OperatorKey ContextKey = "user"

// JWTMiddleware валидирует JWT токен и передает данные пользователя в контекст
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Ожидаем токен в формате: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверка и валидация токена
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.SettingsGlobal.SecretJWT, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Извлечение данных пользователя из токена
		user, err := security.ExtractUserFromClaims(claims)
		if err != nil {
			http.Error(w, "Failed to extract user from token", http.StatusInternalServerError)
			return
		}

		// Передача данных пользователя в контекст
		ctx := context.WithValue(r.Context(), OperatorKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
