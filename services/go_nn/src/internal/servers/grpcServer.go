package servers

import (
	"context"
	"fmt"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
	pb "github.com/PerfectStepCoder/yp_go_nn/src/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"runtime"
	"sync"
)

// ServerGRPC поддерживает все необходимые методы сервера.
type ServerGRPC struct {
	nn *engine.OnnxNeuralNetwork
	pb.UnimplementedClassifyNNServer
	server       *grpc.Server
	countWorkers int
}

// CreateTask - отправляем изображение для классификации.
func (s *ServerGRPC) CreateOneTask(ctx context.Context, in *pb.TaskOneRequest) (*pb.TaskOneResponse, error) {
	var response pb.TaskOneResponse

	image, err := engine.BytesToFloat32Matrix(in.Image, int(in.Height), int(in.Width))

	if err != nil {
		return nil, err
	}

	labelClassNames, err := s.nn.Detect(image)

	response.TaskUID = in.TaskUID
	response.ClassName = labelClassNames[0]

	return &response, nil
}

// CreateBatchTask - отправляем изображение для классификации.
func (s *ServerGRPC) CreateBatchTask(ctx context.Context, in *pb.TaskBatchRequest) (*pb.TaskBatchResponse, error) {
	var response pb.TaskBatchResponse

	images := make([][]float32, len(in.Images))

	for i, image := range in.Images {
		img, err := engine.BytesToFloat32Slice(image)
		if err != nil {
			fmt.Println(err)
		}
		images[i] = img
	}

	labelClassNames, err := s.nn.Detect(images)
	if err != nil {
		fmt.Println(err)
	}

	response.TaskUID = in.TaskUID
	response.ClassNames = labelClassNames

	return &response, nil
}

// CreateBatchCodeTask - отправляем батч изображений для классификации.
func (s *ServerGRPC) CreateBatchCodeTask(ctx context.Context, in *pb.TaskBatchRequest) (*pb.TaskBatchCodeResponse, error) {
	var response pb.TaskBatchCodeResponse

	images := make([][]float32, len(in.Images))

	for i, image := range in.Images {
		img, err := engine.BytesToFloat32Slice(image)
		if err != nil {
			fmt.Println(err)
		}
		images[i] = img
	}

	labelClassCodes, err := s.nn.DetectCode(images)
	if err != nil {
		fmt.Println(err)
	}

	response.TaskUID = in.TaskUID
	response.ClassCodes = engine.IntToInt32Slice(labelClassCodes)

	return &response, nil
}

func (s *ServerGRPC) CreateBatchsTask(ctx context.Context, in *pb.TaskBatchsRequest) (*pb.TaskBatchsResponse, error) {
	var response pb.TaskBatchsResponse

	var wg sync.WaitGroup
	wg.Add(s.countWorkers) // Увеличиваем счетчик для ожидания

	inputCh := make(chan *pb.TaskBatchRequest, len(in.Batchs))
	outputCh := make(chan *pb.TaskBatchResponse, len(in.Batchs))

	for _, batch := range in.Batchs {
		inputCh <- batch
	}

	fmt.Printf("Got batches: %d", len(inputCh))
	close(inputCh)

	// Обработка батчей
	for i := 0; i < s.countWorkers; i++ {
		go func(inputCh chan *pb.TaskBatchRequest, outputCh chan *pb.TaskBatchResponse) {
			defer wg.Done() // уменьшаем счетчик по завершении работы
			for batch := range inputCh {
				result, _ := s.CreateBatchTask(context.Background(), batch)
				outputCh <- result
			}
		}(inputCh, outputCh)
	}

	// Горутинa для закрытия outputCh после завершения всех обработчиков
	go func() {
		wg.Wait()       // Ожидаем завершения всех горутин
		close(outputCh) // Закрываем outputCh, чтобы завершить чтение из него
	}()

	for result := range outputCh {
		response.Batchs = append(response.Batchs, result)
	}

	return &response, nil
}

func NewServerGRPC(nn *engine.OnnxNeuralNetwork) (*ServerGRPC, error) {

	return &ServerGRPC{
		nn:           nn,
		countWorkers: runtime.NumCPU(),
	}, nil
}

func (s *ServerGRPC) Start(addr string) error {

	// Cоздаём gRPC-сервер без зарегистрированной службы
	s.server = grpc.NewServer(
		grpc.MaxRecvMsgSize(50*1024*1024), // Максимальный размер принимаемого сообщения — 50 МБ
		grpc.MaxSendMsgSize(50*1024*1024), // Максимальный размер отправляемого сообщения — 50 МБ
	)

	// Регистрируем сервис
	pb.RegisterClassifyNNServer(s.server, s)

	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln(err)
	}

	return s.server.Serve(listen)
}

func (s *ServerGRPC) Stop(ctx context.Context) error {
	// Остановка сервера с плавным завершением
	s.server.GracefulStop()
	return nil
}

func (s *ServerGRPC) GetInfo(ctx context.Context, in *emptypb.Empty) (*pb.ServiceInfoNN, error) {
	var response pb.ServiceInfoNN
	response.Name = "YoloNN"
	response.Description = "[Go] Trained Fashion MNIST"
	response.Version = "1.0.0.0"
	return &response, nil
}
