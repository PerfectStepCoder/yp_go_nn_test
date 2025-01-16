package main

import (
	"bufio"
	"fmt"
	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"os"
)

func startReporters(connectorOne pb.ClassifyNNClient, connectorTwo pb.ClassifyNNClient) {

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("Введите номер команды:\n",
			"\t0-выход\n",
			"\t1-сравнение идентичности работы двух сервисов\n",
			"\t2-работа сервиса для одного изображения\n",
			"\t3-производительность сервиса\n",
		)

		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if input == "0" {
			break
		}

		switch input {
		case "1":
			DoReportCompareBatch(connectorOne, connectorTwo)
			continue
		case "2":
			DoReportSolo(connectorOne)
			continue
		case "3":
			DoReportPerformance(connectorTwo)
			continue
		default:
			fmt.Println("Команда не распознана")
			continue
		}
	}
	fmt.Println("Завершение работы утилиты")
}
