package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const (
	// SQL-схема для создания таблицы
	schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL
);

CREATE INDEX id_date ON scheduler(date);
`
	//Имя БД
	Name = "scheduler.db"
)

// Глобальная переменная для доступа к БД
var (
	DB *sql.DB
)

func Init(dbFile string) error {

	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем соединение с БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("can't open DB: %w", err)
	}

	// Создание таблицы, если ее нет
	if install {
		if _, err = DB.Exec(schema); err != nil {
			return fmt.Errorf("can't create DB: %w", err)
		}
		fmt.Println("DB created")
	}

	return nil
}
