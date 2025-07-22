package main

import (
	"Final_project/db"
	"Final_project/tests"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

func main() {

	if err := db.Init(db.Name); err != nil {
		fmt.Println(err)
		return
	}
	defer db.DB.Close()

	//Директория для возвращаемых файлов
	webDir := "./web"

	//Порт из переменной пакета tests
	port := strconv.Itoa(tests.Port)

	// Настройка хэндлера
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Запуск сервера
	fmt.Printf("server started at port: %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("server start error: %v\n", err)
	}
}
