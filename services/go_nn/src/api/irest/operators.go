package httpp // http protocol

import (
	"encoding/json"
	"net/http"

	//"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/security"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	"github.com/go-chi/chi/v5"
)

// @Summary      Get operator by UID
// @Description  Get operator
// @Tags         Operators
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        operatorUID   path      string  true  "Operator UID"
// @Success      200   body  models.Operator true "Operator"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      404  {string}  string  "Operator not found"
// @Router       /operators/id/{operatorUID} [get]
func GetOperatorByUIDHandler(mainStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Извлечение пользователя из контекста
		_, ok := req.Context().Value(OperatorKey).(*security.OperatorJWT)
		if !ok {
			http.Error(res, "Operator not found in context", http.StatusUnauthorized)
			return
		}

		// Извлечение параметра operatorUID из URL
		operatorUID := chi.URLParam(req, "operatorUID")

		if operatorUID == "" {
			http.Error(res, "operatorUID is required", http.StatusBadRequest)
			return
		}

		result, err := mainStorage.GetOperatorByUID(operatorUID)
		if err != nil {
			http.Error(res, "Operator not found", http.StatusNotFound)
			return
		}
		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(result)
		if err != nil {
			http.Error(res, "Error writing response", http.StatusInternalServerError)
			return
		}

		res.Write(jsonResp)

		res.WriteHeader(http.StatusOK)

	}
}

// @Summary      Get operator by Name
// @Description  Get operator
// @Tags         Operators
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name   path      string  true  "Operator name"
// @Success      200   body  models.Operator true "Operator"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      404  {string}  string  "Operator not found"
// @Router       /operators/name/{name} [get]
func GetOperatorByNameHandler(mainStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Извлечение пользователя из контекста
		_, ok := req.Context().Value(OperatorKey).(*security.OperatorJWT)
		if !ok {
			http.Error(res, "Operator not found in context", http.StatusUnauthorized)
			return
		}

		// Извлечение параметра name из URL
		name := chi.URLParam(req, "name")

		if name == "" {
			http.Error(res, "Name is required", http.StatusBadRequest)
			return
		}

		result, err := mainStorage.GetOperatorByName(name)
		if err != nil {
			http.Error(res, "Operator not found", http.StatusNotFound)
			return
		}
		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(result)
		if err != nil {
			http.Error(res, "Error writing response", http.StatusInternalServerError)
			return
		}

		res.Write(jsonResp)

		res.WriteHeader(http.StatusOK)

	}
}
