package main

import (
	"fmt"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
)

func main() {

	inputLayer := engine.NeuralLayer{
		Name:  "input",
		Shape: []int64{1, 28, 28},
	}
	outputLayer := engine.NeuralLayer{
		Name:  "output",
		Shape: []int64{10},
	}

	nn := engine.NewOnnxNeuralNetwork("../models/yolo_fashion_mnist.onnx",
		"../lib/libonnxruntime.1.20.1.dylib", inputLayer, outputLayer)

	// Загрузка датасета
	images, labels, err1 := engine.LoadDataset("../../datasets/fashion_mnist_test.csv", 105) // labels

	if err1 != nil {
		fmt.Printf("Ошибка загрузки датасета: %v\n", err1)
		return
	}

	for i, image_batch := range images {

		classCode, _ := nn.DetectCode(image_batch)

		fmt.Println(engine.CalculateMatchPercentage(labels[i], classCode))

		// if engine.CompareArrays(labels[i], classCode) == false {
		// 	fmt.Println("Not match!")
		// }
		// break
		//fmt.Println(nn.DetectCode(image_batch))
	}
	fmt.Println("Done")
}
