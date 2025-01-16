package main

import (
	"context"
	"fmt"

	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DoReportSolo - отправляет одно изображение одной нейронной сети
func DoReportSolo(c pb.ClassifyNNClient) {

	// Загрузка датасета
	batchSize := 200
	images, _, err := engine.LoadDataset("../../../../../data/datasets/fashion_mnist_test.csv", batchSize) // labels

	if err != nil {
		fmt.Printf("Error load of dataset: %v\n", err)
		return
	}

	for i, imageBatch := range images {

		imageBatchBytes, err := engine.Float32MatrixToBytes(engine.GetFirstImage(imageBatch))
		if err != nil {
			fmt.Printf("Error Float32MatrixToBytes batch: %v\n", err)
			return
		}

		//height, width := engine.GetMatrixSize(imageBatch)
		height, width := 1, len(imageBatch[0])

		requestTaskOne := pb.TaskOneRequest{
			TaskUID: uuid.New().String(),
			Image:   imageBatchBytes,
			Width:   int32(width),
			Height:  int32(height),
		}

		result, err := c.CreateOneTask(context.Background(), &requestTaskOne)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(i, result.ClassName)

	}

}

// DoReportBatch - сравнение идентичности работы двух сервисов
func DoReportCompareBatch(connectorOne pb.ClassifyNNClient, connectorTwo pb.ClassifyNNClient) {

	// Общая информация об нейронных сетях
	infoOne, _ := connectorOne.GetInfo(context.Background(), &emptypb.Empty{})
	infoTwo, _ := connectorTwo.GetInfo(context.Background(), &emptypb.Empty{})

	fmt.Printf("NN one:\nName: %s\nDescription: %s\nVersion: %s\n", infoOne.Name, infoOne.Description, infoOne.Version)
	fmt.Printf("NN one:\nName: %s\nDescription: %s\nVersion: %s\n", infoTwo.Name, infoTwo.Description, infoTwo.Version)

	// Загрузка датасета
	batchSize := 200
	images, labels, err := engine.LoadDataset("../../../../../data/datasets/fashion_mnist_test.csv", batchSize) // labels

	if err != nil {
		fmt.Printf("Error load of dataset: %v\n", err)
		return
	}

	var totalAccuracyOne, totalAccuracyTwo float64

	for i, imageBatch := range images {

		imageBatchBytes, err := engine.Float32ToBytes2D(imageBatch)
		if err != nil {
			fmt.Printf("Error Float32MatrixToBytes batch: %v\n", err)
			return
		}

		//height, width := engine.GetMatrixSize(imageBatch)

		requestTaskBatch := pb.TaskBatchRequest{
			TaskUID: uuid.New().String(),
			Images:  imageBatchBytes,
			Width:   int32(28), // размер изображения
			Height:  int32(28), // размер изображения
		}

		resultOne, errOne := connectorOne.CreateBatchCodeTask(context.Background(), &requestTaskBatch) //CreateBatchTask
		if errOne != nil {
			fmt.Println(errOne)
		}

		resultTwo, errTwo := connectorTwo.CreateBatchCodeTask(context.Background(), &requestTaskBatch) //CreateBatchTask
		if errTwo != nil {
			fmt.Println(errTwo)
		}

		if engine.CompareArrays(resultOne.ClassCodes, resultTwo.ClassCodes) {
			fmt.Printf("Batch: %d Total match\n", i)
		} else {
			rateMatch := engine.CalculateMatchPercentage(resultOne.ClassCodes, resultTwo.ClassCodes)
			fmt.Printf("Batch: %d Only match: %f\n", i, rateMatch)
		}

		labels_batch := engine.IntToInt32Slice(labels[i])
		totalAccuracyOne = totalAccuracyOne + engine.CalculateMatchPercentage(resultOne.ClassCodes, labels_batch)
		totalAccuracyTwo = totalAccuracyTwo + engine.CalculateMatchPercentage(resultTwo.ClassCodes, labels_batch)

	}

	fmt.Println("Avg accuracy one: ", totalAccuracyOne/float64(len(images)))
	fmt.Println("Avg accuracy two: ", totalAccuracyTwo/float64(len(images)))

}

// DoReportPerformance - производительность сервиса
func DoReportPerformance(c pb.ClassifyNNClient) {
	// Загрузка датасета
	batchSize := 200
	images, _, err := engine.LoadDataset("../../../../../data/datasets/fashion_mnist_test.csv", batchSize) // labels

	if err != nil {
		fmt.Printf("Error load of dataset: %v\n", err)
		return
	}

	var batchsRequest pb.TaskBatchsRequest

	for _, imageBatch := range images {

		imageBatchBytes, err := engine.Float32ToBytes2D(imageBatch)
		if err != nil {
			fmt.Printf("Error Float32MatrixToBytes batch: %v\n", err)
			return
		}

		//height, width := engine.GetMatrixSize(imageBatch)

		requestTaskBatch := pb.TaskBatchRequest{
			TaskUID: uuid.New().String(),
			Images:  imageBatchBytes,
			Width:   int32(28), // размер изображения
			Height:  int32(28), // размер изображения
		}
		batchsRequest.Batchs = append(batchsRequest.Batchs, &requestTaskBatch)
	}

	result, err := c.CreateBatchsTask(context.Background(), &batchsRequest)
	if err != nil {
		fmt.Printf("CreateBatchsTask: %v\n", err)
		return
	}

	for i, batch := range result.Batchs {
		fmt.Println(i, batch.ClassNames)
	}
}
