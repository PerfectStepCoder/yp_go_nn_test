package main

import (
	"fmt"
	"log"

	// "context"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	// "github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	// "github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"

	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var setting, _ = NewSettings()

func main() {

	fmt.Println(setting)

	// Устанавливаем соединение с сервером
	connOne, errOne := grpc.NewClient(fmt.Sprintf("%s:%s", setting.ServiceHostOne, setting.ServicePortOne),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errOne != nil {
		log.Fatal(errOne)
	}
	defer connOne.Close()

	connTwo, errTwo := grpc.NewClient(fmt.Sprintf("%s:%s", setting.ServiceHostTwo, setting.ServicePortTwo),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errTwo != nil {
		log.Fatal(errTwo)
	}
	defer connTwo.Close()

	// Получаем переменную интерфейсного типа ClassifyNNClient, через которую будем отправлять сообщения
	connectorOne := pb.NewClassifyNNClient(connOne)
	connectorTwo := pb.NewClassifyNNClient(connTwo)

	// Функция для выполнения тестирования
	startReporters(connectorOne, connectorTwo)

}
