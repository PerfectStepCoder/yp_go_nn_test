package httpp // http protocol

import (
	"encoding/json"
	"net/http"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/models"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/security"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	//"github.com/go-chi/chi/v5"
)

// @Summary      Task
// @Description  Create Task with one image
// @Tags         Tasks
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file formData file true "File image"
// @Param        taskUID formData string true "UUID of the task"
// @Success      200   body  models.TaskOneResponse true "Task one image"
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /tasks/one  [post]
func TaskOneHandler(mainStorage storage.Storage, nn *engine.OnnxNeuralNetwork) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Извлечение пользователя из контекста
		_, ok := req.Context().Value(OperatorKey).(*security.OperatorJWT)
		if !ok {
			http.Error(res, "Operator not found in context", http.StatusUnauthorized)
			return
		}

		// Попытка извлечь файл из запроса
		file, _, err := req.FormFile("file") // fileHeader
		if err != nil {
			http.Error(res, "Did not catch the file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close() // Закрыть файл после обработки

		taskUID := req.FormValue("taskUID")
		if taskUID == "" {
			http.Error(res, "Missing taskUID parameter", http.StatusBadRequest)
			return
		}

		// Вывод информации о загруженном файле
		//fmt.Fprintf(res, "Загружен файл: %s (размер: %d байт)\n", fileHeader.Filename, fileHeader.Size)

		image, err := ImageToFloat32Matrix(file)

		if err != nil {
			http.Error(res, "Did not worked ImageToFloat32Matrix", http.StatusInternalServerError)
			return
		}

		var lblClassName string
		if labelClassName, err := nn.Detect(ReshapeTo1xN(image)); err != nil {
			http.Error(res, "Did not worked Neural Network", http.StatusInternalServerError)
			return
		} else {
			lblClassName = labelClassName[0]
		}

		result := models.TaskOneResponse{
			TaskUID: taskUID,
			Label:   lblClassName,
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
