package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Tasks(i int, DB *sql.DB) ([]*TaskDTO, error) {

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", i)
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
	rows, err := DB.Query("SELECT * FROM scheduler WHERE date = ? LIMIT ?", date, i)

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
	rows, err := DB.Query("SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search), i)

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
