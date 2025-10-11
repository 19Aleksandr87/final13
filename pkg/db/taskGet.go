package db

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"
)

func UpdateTask(s string) (*TaskDTO, error) {
	DB, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	var task Task
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", i)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	dto := convertTaskToTaskDTO(task)
	return &dto, nil
}
