package db

import (
	"errors"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	if DB == nil {
		return 0, errors.New("database not initialized")
	}

	/* query := `
	    INSERT INTO scheduler (date, title, comment, repeat)
	    VALUES (:date, :title, :comment, :repeat)
	    `
		res, err := DB.Exec(query,
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
		)*/

	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
