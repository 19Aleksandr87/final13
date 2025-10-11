package db

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"
)

func convertTaskToTaskDTO(task Task) TaskDTO {
	return TaskDTO{
		ID:      strconv.Itoa(task.ID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}

func Tasks(i int) ([]*TaskDTO, error) {
	DB, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", i)
	if err != nil {
		return nil, err
	}

	var tasks []*TaskDTO
	for rows.Next() {
		var task Task
		var t TaskDTO
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		t = convertTaskToTaskDTO(task)

		tasks = append(tasks, &t)
	}
	return tasks, nil
}
