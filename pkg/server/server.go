package server

import (
	"Final_project/pkg/api"
	"Final_project/tests"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	logger *log.Logger
	port   string
}

func NewServer(logger *log.Logger) *Server {

	//Порт из переменной пакета tests
	port := strconv.Itoa(tests.Port)

	return &Server{
		logger: logger,
		port:   port,
	}
}

func (s *Server) Start() error {

	api.Init()

	//Директория для возвращаемых файлов
	webDir := "./web"

	// Настройка хэндлера
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Запуск сервера
	s.logger.Println("The server is running on http://localhost:" + s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
