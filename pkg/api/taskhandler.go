package api

import (
	"Final_project/pkg/db"
	"encoding/json"
	"net/http"
	"time"
)

var TaskLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Общий хэнделр для единичной задачи, разбивающий по методу обращения
func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Хэндлер для списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(TaskLimit)
	if err != nil {
		writeJSONerror(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TasksResp{Tasks: tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Получение задачи
func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONerror(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONerror(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

// Изменение задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	//Декодирование
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONerror(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// Проверка наличия заголовка
	if task.ID == "" {
		writeJSONerror(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Проверка наличия заголовка
	if task.Title == "" {
		writeJSONerror(w, "task title required", http.StatusBadRequest)
		return
	}

	// проверка даты
	if err := checkDate(&task); err != nil {
		writeJSONerror(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {

		writeJSONerror(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}

// Отметить задачу сделанной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONerror(w, "ID is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONerror(w, "Task not found", http.StatusBadRequest)
		return
	}

	// Если нет правил повторения, удаляем таску
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONerror(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)

	next, err := NextDate(now, task.Date, task.Repeat)

	if err != nil {
		writeJSONerror(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		writeJSONerror(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})

}

// Удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONerror(w, "No task with this ID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSONerror(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}
