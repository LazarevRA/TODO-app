package main

import (
	"Final_project/pkg/db"
	"Final_project/pkg/server"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {

	// Подключение к БД
	if err := db.Init(db.Name); err != nil {
		fmt.Println(err)
		return
	}
	//Не забываем закрыть БД
	defer db.DB.Close()

	//Новый логгер и сервер
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	srv := server.NewServer(logger)

	//Запуск сервера, api.Init встроен в server.Start
	if err := srv.Start(); err != nil {
		logger.Fatal("Error starting server: ", err)
	}

}
