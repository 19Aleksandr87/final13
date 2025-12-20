package db

import (
	"database/sql"
	"fmt"
)

func Tasks(i int, DB *sql.DB) ([]*TaskDTO, error) {

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1", i)
	if err != nil {
		return nil, err
	}

	var tasks []*TaskDTO
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		t := convertTaskToTaskDTO(task)

		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func DateTasks(date string, i int, DB *sql.DB) ([]*TaskDTO, error) {
	rows, err := DB.Query("SELECT * FROM scheduler WHERE date = $1 LIMIT $2", date, i)

	if err != nil {
		return nil, err
	}

	var tasks []*TaskDTO

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		dto := convertTaskToTaskDTO(task)
		tasks = append(tasks, &dto)
	}

	return tasks, nil
}

func SearchTasks(search string, i int, DB *sql.DB) ([]*TaskDTO, error) {
	rows, err := DB.Query("SELECT * FROM scheduler WHERE title LIKE $1 OR comment LIKE $2 ORDER BY date LIMIT $3", fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search), i)

	if err != nil {
		return nil, err
	}

	var tasks []*TaskDTO

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		dto := convertTaskToTaskDTO(task)
		tasks = append(tasks, &dto)
	}

	return tasks, nil
}
