package servers

import (
	"context"
	httpp "github.com/PerfectStepCoder/yp_go_nn/src/api/irest"
	_ "github.com/PerfectStepCoder/yp_go_nn/src/api/irest/docs" // Импорт сгенерированной документации
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/engine"
	"github.com/PerfectStepCoder/yp_go_nn/src/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type ServerHTTP struct {
	storage storage.Storage
	nn      *engine.OnnxNeuralNetwork
	router  *chi.Mux
	server  *http.Server
}

func NewHTTPServer(mainStorage storage.Storage, nn *engine.OnnxNeuralNetwork) (*ServerHTTP, error) {
	return &ServerHTTP{
		storage: mainStorage,
		nn:      nn,
		router:  chi.NewRouter(),
	}, nil
}

func (s *ServerHTTP) Start(addr string) error {

	// Swagger endpoint
	s.router.Mount("/swagger", httpSwagger.WrapHandler)

	// Register routes
	s.router.Post("/register", httpp.RegisterHandler(s.storage))
	s.router.Post("/login", httpp.LoginHandler(s.storage))

	// Защищенные маршруты
	s.router.Group(func(r chi.Router) {
		r.Use(httpp.JWTMiddleware)
		// Operators
		r.Route("/operators", func(r chi.Router) {
			r.Get("/id/{operatorUID}", httpp.GetOperatorByUIDHandler(s.storage))
			r.Get("/name/{name}", httpp.GetOperatorByNameHandler(s.storage))
		})
		// Tasks
		r.Route("/tasks", func(r chi.Router) {
			r.Post("/one", httpp.TaskOneHandler(s.storage, s.nn))
		})
	})

	// Инициализируем HTTP сервер
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}

func (s *ServerHTTP) Stop(ctx context.Context) error {
	// Остановка сервера с плавным завершением
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
