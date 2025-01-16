package main

import (
	"context"
	"fmt"

	//"log"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"github.com/google/uuid"
)

func DoReport(c pb.ClassifyNNClient) {

	// Загрузка датасета
	batchSize := 200
	images, labels, err := engine.LoadDataset("../../../../../data/datasets/fashion_mnist_test.csv", batchSize)  // labels

	if err != nil {
		fmt.Printf("Error load of dataset: %v\n", err)
		return
	}

	for i, imageBatch := range images {

		imageBatchBytes, err := engine.Float32MatrixToBytes(imageBatch)
		if err != nil {
			fmt.Printf("Error Float32MatrixToBytes batch: %v\n", err)
			return
		}

		height, width := engine.GetMatrixSize(imageBatch)

		requestTaskOne := pb.TaskOneRequest{
			TaskUID: uuid.New().String(),
			Image: imageBatchBytes,
			Width: int32(width),
			Height: int32(height),
			
		}

		result, err := c.CreateOneTask(context.Background(), &requestTaskOne)
		
		fmt.Println(i, result.ClassName)
		
	}

}