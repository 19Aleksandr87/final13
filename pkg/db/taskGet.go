package db

import (
	"database/sql"
	"strconv"
)

func UpdateTask(s string, DB *sql.DB) (*TaskDTO, error) {

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	var task Task
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1", i)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	dto := convertTaskToTaskDTO(task)
	return &dto, nil
}
