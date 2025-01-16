# Комплекс для запуска нейронных сетей

<details>

<summary>Описание сервиса и его базовой функциональности</summary>

### Цели и задачи проекта

Основная цель проекта разработать механизм переноса состояние нейронной сети из окружения на Python в окружение на Go

Задачи решаемые в данном проекте:
- разработка функциональности загрузка модели нейросети (в формате onnx), в сервис на Go для последующей работы с данной моделью для решения задач классификации изображений;
- разработка утилиты для сравнительного анализа сервисов на Python и Go которые содержат одинаковые модели нейронной сети.

</details>

<details>

<summary>Кодовая база проекта</summary>

### Статический анализ кода

#### Проверка стиля 

```Go
gofmt -s -w .
```
-s simplifies the code\
-w writes results directly

```Go
gofmt -l -s .
goimports -l .
```
/Users/dmitrii/go/bin/goimports -l . (локальный запуск пример)

#### Проверка стуктурных тегов
```Go
go vet -structtag
```

#### Загрузка зависимостей
```Go
cd ./services/go_nn/src
go mod tidy
```

</details>

<details>

<summary>Описание файловой структуры проекта</summary>

# Файловая структура проекта
| Путь    | Назначение |
| -------- | ------- |
| api  | Реализация REST API      |
| api/docs | Документация Swagger     |
| cmd    | Точка запуска проекта    |
| cmd/main.go    | Запуск сервиса (rest или grpc)   |
| cmd/console    | Отладка и тесты   |
| cmd/test_reporter    | Утилита сранение нейросетевых сервисов|
| configs   | Настройки в проекте   |
| internal   | Внутренее пакеты проекта   |
| internal/engine   | Функционал для работы с нейронной сетью|
| internal/engine/loader   | Загрузка данных тестового набора, функции по преобразованию массивов|
| internal/engine/mnist   | Работа с классами Fashion MNIST|
| internal/engine/nn   | Нейронная сеть и оnnx формат|
| internal/engine/utils   | Функции по работе с массивами|
| models   | Модели сущностей: операторы (пользователи и инженеры) и задачи   |
| proto   | Описание grpc сущностей а также паки с результатами кодогенерации |
| security   | Работа с JWT и паролями для REST API|
| servers   | Реализация нейросетевых сервисов|
| servers/httpServer   | REST API сервис   |
| servers/grpcServer   | GRPC сервис  |
| storage   | Модели для хранения в базе данных   |
| storage/database   | Реализация хранилища на Postgres |

</details>

<details>

<summary>Переменные окружения</summary>

### Хранилище
POSTGRES_USER=admin\
POSTGRES_PASSWORD=password\
POSTGRES_DB=nn\
POSTGRES_HOST=localhost\
POSTGRES_PORT=6443\
### Сервис
SERVICE_HOST=localhost\
SERVICE_PORT=8011\
SERVICE_PROTOCOL=http\
SECRET_JWT=ffduyfdb3534Gfdituetr\
### Для утилиты с отчетами
SERVICE_HOST_ONE=localhost\
SERVICE_PORT_ONE=50051\
SERVICE_HOST_TWO=localhost\
SERVICE_PORT_TWO=3001\

</details>

<details>

<summary>Описание интерфейсов всех сервисов (API)</summary>

### REST HTTP
Создание документации Swagger
```Go
go install github.com/swaggo/swag/cmd/swag@latest
export PATH="$PATH:$(go env GOPATH)/bin"  // чтобы видна была утилита
swag init -g cmd/main.go -o api/irest/docs // создаем документацию
```

### GRPC
#### Кодогенерация
##### Go доступность в консоле
```Go
export PATH="$PATH:$(go env GOPATH)/bin"
cd ../internal/proto
protoc --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative models.proto server.proto
```
##### Python
Установка
```Python
pip install --upgrade grpcio grpcio-tools
python -m grpc_tools.protoc --version
```
Генерация
```Python
python -m grpc_tools.protoc -I . --python_out=./gen_py --pyi_out=./gen_py --grpc_python_out=./gen_py ./models.proto ./server.proto
```
-I .: Указывает путь к директории, где находятся файлы .proto. Текущая директория обозначена как .
--python_out=.: Указывает путь для генерации файлов Python, содержащих структуры данных из .proto
--grpc_python_out=.: Указывает путь для генерации файлов Python, содержащих определения gRPC-сервисов 
example.proto: Имя файла .proto, который вы хотите обработать.

</details>


<details>

<summary>Запуск сервисов</summary>

Внимание!!!
Для запуска нужны файлы: дамп нейросети, библиотека распакомки состояний нейросети в формате onnx, тестовый датасет. Всего этого нет в репозитории. Обратитесь к разработчику.

#### Запуск сервиса (Python)
Установка ultralytics
```Python
conda install -c conda-forge ultralytics
```
#### Запуск сервиса (Go)
```Go
cd ../cmd/
go run main.go -h 0.0.0.0 -p 3001 -m grpc
```
#### Запуск утилиты создание отчетов
```Go
cd ../cmd/test_reporter
go run main.go reporter.go settings.go menu.go
```

</details>

<details>

<summary>Запуск тестов</summary>

#### Запуск тестов для сервиса на Go

Внимание!!!
Для выполнения тестов нужны файлы: дамп нейросети, библиотека распакомки состояний нейросети в формате onnx, тестовый датасет. Всего этого нет в репозитории, но я локально выполнил тесты и разместил результаты покрытия тестами в файлы:
- coverage.out
- coverage.html (этот файл можно открыть в браузере для просмотра)

```Go
go test ./...
go test ./services/go_nn/src/... -coverprofile=coverage.out
```
### Выводим отчет о покрытие тестами кода
```Go
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

</details>
