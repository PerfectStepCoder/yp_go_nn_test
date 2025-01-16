package main

import (
	"fmt"
	"log"
	// "context"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	"github.com/PerfectStepCoder/yp_go_nn/src/configs"
	// "github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	// "github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"

	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var settings = configs.SettingsGlobal

func main() {

	configs.ParseFlags(settings)
	fmt.Println(settings)

	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", settings.ServiceHost, settings.ServicePort), 
	                            grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Получаем переменную интерфейсного типа ClassifyNNClient, через которую будем отправлять сообщения
	c := pb.NewClassifyNNClient(conn)

	// Функция для выполнения тестирования
	DoReport(c)
}