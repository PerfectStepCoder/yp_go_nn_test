package configs

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Функция для создания логгера
func GetLogger() *logrus.Logger {

	logger := logrus.New()

	// Устанавливаем вывод логгера на несколько writer'ов (консоль и файл)
	multiWriter := io.MultiWriter(os.Stdout)
	logger.SetOutput(multiWriter)

	// Устанавливаем формат вывода
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}
