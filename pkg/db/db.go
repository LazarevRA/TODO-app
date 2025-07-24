package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const (
	// SQL-схема для создания таблицы
	schema = "CREATE TABLE scheduler (\nid INTEGER PRIMARY KEY AUTOINCREMENT,\ndate CHAR(8) NOT NULL DEFAULT '',\ntitle VARCHAR(255) NOT NULL DEFAULT '',\ncomment TEXT NOT NULL DEFAULT '',\nrepeat VARCHAR(128) NOT NULL DEFAULT ''\n);"

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

		if _, err = DB.Exec("CREATE INDEX idx_date ON scheduler(date);"); err != nil {
			return fmt.Errorf("can't create index on date: %w", err)
		}
		fmt.Println("date index created")
	}

	return nil
}
