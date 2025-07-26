package api

import (
	"Final_project/pkg/db"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ok
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task
	var err error
	//Декодирование
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, "error", fmt.Sprintf("incorrect data in JSON: %s", err))

		return
	}

	//Проверка на наличие заголовка
	if task.Title == "" {

		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, "error", "title required")
		return
	}

	//Проверка даты
	if err := checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, "error", fmt.Sprintf("Check date error: %s", err))
		return
	}

	//Добавляем задачу в БД
	id, err := db.AddTask(&task)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, "error", fmt.Sprintf("DB add task error: + %s", err))
		return
	}
	//Запись last ID

	writeJSON(w, "id", id)

}

// OK
func writeJSON(w http.ResponseWriter, key string, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]any{key: value})
}

// OK
func checkDate(task *db.Task) error {

	now := time.Now()
	now = now.Truncate(24 * time.Hour)

	// Если дата не указана, устанавливаем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(Layout)
	}
	// парсим дату
	date, err := time.Parse(Layout, task.Date)
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
	if afterNow(now, date) {
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			task.Date = next
		}
	}

	return nil
}
