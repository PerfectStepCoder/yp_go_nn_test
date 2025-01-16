package httpp // http protocol

import (
	"encoding/json"
	"net/http"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/security"
)

// RegisterHandler handles operator registration
// @Summary      Register a new operator
// @Description  Register a new operator with name and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.OperatorRequest true  "Operator data" example({"name": "john_doe", "password": "securepassword"})
// @Success      201   body  models.UserResponseUID true "Operator created"
// @Failure      400  {string}  string  "Invalid input"
// @Failure      409  {string}  string  "Operator already exists"
// @Router       /register [post]
func RegisterHandler(mainStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var operator models.Operator
		if err := json.NewDecoder(req.Body).Decode(&operator); err != nil {
			http.Error(res, "Invalid input", http.StatusBadRequest)
			return
		}

		if hashed, hashErr := security.HashPassword(operator.Password); hashErr == nil {
			operator.Password = hashed
		}

		if _, err := mainStorage.CreateOperator(&operator); err != nil {
			if err.Error() == "UNIQUE constraint failed: users.username, users.email" {
				http.Error(res, "User already exists", http.StatusConflict)
			} else {
				http.Error(res, "Failed to create user", http.StatusInternalServerError)
			}
			return
		}

		result := models.OperatorResponseUID{
			OperatorUID: operator.OperatorUID.String(),
		}

		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(result)
		if err != nil {
			http.Error(res, "Error writing response", http.StatusInternalServerError)
			return
		}
		res.Write(jsonResp)

		res.WriteHeader(http.StatusCreated)
	}
}

// LoginHandler аутентифицирует пользователя и возвращает JWT токен
// @Summary      Login a operator
// @Description  Authenticate operator and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login  body  models.LoginRequest  true  "Login credentials"
// @Success      200  {object}  models.LoginResponse  "JWT Token"
// @Failure      400  {string}  string  "Invalid input"
// @Failure      401  {string}  string  "Invalid credentials"
// @Router       /login [post]
func LoginHandler(mainStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Получаем пользователя из базы данных
		operator, err := mainStorage.GetOperatorByName(req.Name)
		if err != nil || !security.CheckPassword(operator.Password, req.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Генерация JWT токена
		token, err := security.GenerateJWT(operator)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// Возвращаем токен в ответе
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
		w.WriteHeader(http.StatusOK)
	}
}
