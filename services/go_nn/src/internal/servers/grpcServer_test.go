package servers

import (
	"context"
	//"fmt"
	"testing"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDoReportSolo(t *testing.T) {
	
	// NeuralNetwork
	inputLayer := engine.NeuralLayer{
		Name:  "input",
		Shape: []int64{1, 28, 28},
	}
	outputLayer := engine.NeuralLayer{
		Name:  "output",
		Shape: []int64{10},
	}

	nn := engine.NewOnnxNeuralNetwork("/Users/dmitrii/EducationProjects/yp_go_nn_test/services/go_nn/src/models/yolo_fashion_mnist.onnx",
		"/Users/dmitrii/EducationProjects/yp_go_nn_test/services/go_nn/src/lib/libonnxruntime.1.20.1.dylib", inputLayer, outputLayer)

	// Создаём экземпляр сервиса
	server, _ := NewServerGRPC(nn)

	// Загрузка датасета
	batchSize := 200
	images, labels, _ := engine.LoadDataset("/Users/dmitrii/EducationProjects/yp_go_nn_test/data/datasets/fashion_mnist_test.csv", batchSize) // labels

	for i, imageBatch := range images {

		imageBatchBytes, err := engine.Float32ToBytes2D(imageBatch)
		assert.NoError(t, err)

		requestTaskBatch := pb.TaskBatchRequest{
			TaskUID: uuid.New().String(),
			Images:  imageBatchBytes,
			Width:   int32(28), // размер изображения
			Height:  int32(28), // размер изображения
		}

		result, err := server.CreateBatchCodeTask(context.Background(), &requestTaskBatch)
		assert.NoError(t, err)
		
		resultCompare := engine.CalculateMatchPercentage(engine.IntToInt32Slice(labels[i]), result.ClassCodes)
		assert.True(t, resultCompare > 80.0, resultCompare)
	}

}