package main

import (
	"Final_project/tests"
	"fmt"
	"net/http"
	"strconv"
)

func main() {

	//Директория для возвращаемых файлов
	webDir := "./web"

	//Порт из переменной пакета tests
	port := strconv.Itoa(tests.Port)

	// Настройка хэндлера
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Запуск сервера
	fmt.Printf("Сервер запущен на порту %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
