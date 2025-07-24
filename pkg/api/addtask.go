package api

import (
	"Final_project/pkg/db"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task
	//Декодирование
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {

		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	//Проверка на наличие заголовка
	if task.Title == "" {

		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	//Проверка даты
	if err := checkDate(&task); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Добавляем задачу в БД
	id, err := db.AddTask(&task)

	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	//Запись last ID

	writeJSON(w, map[string]int64{"id": id})

}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {

	now := time.Now()

	// Если дата не указана, устанавливаем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	// парсим дату
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string

	// Если правило повторения указано, вычисляем следующую дату
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Проверка, что дата больше чем сейчас
	if afterNow(date) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}
