package db

import (
	"database/sql"
	"fmt"
)

func TaskPut(dto *TaskDTO, DB *sql.DB) error {

	task, err := ConvertToTaskDTO_Task(*dto)
	if err != nil {
		return err
	}
	res, err := DB.Exec("UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}
