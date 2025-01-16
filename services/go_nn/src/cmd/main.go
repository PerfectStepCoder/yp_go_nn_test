// @title           Нейросети API
// @version         1.0
// @description     Нейронные сети по детекции объектов на изображениях
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PerfectStepCoder/yp_go_nn/src/configs"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/servers"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
)

var settings = configs.SettingsGlobal

func main() {
	
	configs.ParseFlags(settings)
	fmt.Println(settings)

	// Storage
	mainStorage, errStorage := storage.NewStorageInPostgres(settings.DatabaseDSN)
	if errStorage != nil {
		log.Printf("data base doesnt work! Error: %s", errStorage)
		return
	}
	defer mainStorage.Close()

	// NeuralNetwork
	inputLayer := engine.NeuralLayer{
		Name: "input",
		Shape: []int64{1, 28, 28},
	}
	outputLayer := engine.NeuralLayer{
		Name: "output",
		Shape: []int64{10},
	}

	nn := engine.NewOnnxNeuralNetwork("../models/yolo_fashion_mnist.onnx", 
									  "../lib/libonnxruntime.1.20.1.dylib", inputLayer, outputLayer)
	
	// Канал для получения системных сигналов (например, SIGINT, SIGTERM)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	var server servers.Server = nil
	var serverErr error

	switch settings.ServiceProtocol {
		case "http":
			server, serverErr = servers.NewHTTPServer(mainStorage, nn)
		case "grpc":
			server, serverErr = servers.NewServerGRPC(nn)
		default:
			log.Fatalf("Not support protocol: %s", settings.ServiceProtocol)
	}

	if serverErr != nil {
		log.Fatalf("Failed to create server: %v", serverErr)
	}	
	
	log.Printf("Service is starting host: %s, port: %s", settings.ServiceHost, settings.ServicePort)
	serverErr = server.Start(fmt.Sprintf("%s:%s", settings.ServiceHost, settings.ServicePort))

	if serverErr != nil {
		log.Fatalf("Failed to start server: %v", serverErr)
	}

	// Блокировка до получения сигнала
	sig := <-signalChan
	log.Printf("Singal: %v. Shutdown server ...\n", sig)
	// Контекст с таймаутом для graceful shutdown (5 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем сервер, даем время для завершения активных запросов
	if err := server.Stop(ctx); err != nil {
		configs.Logger.Fatalf("failed to stop server: %v\n", err)
	}

	log.Printf("Service shutdown successfully")
}