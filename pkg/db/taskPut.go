package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func TaskPut(dto *TaskDTO, DB *sql.DB) error {

	task, err := ConvertToTaskDTO_Task(*dto)
	if err != nil {
		return err
	}
	res, err := DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.ID))
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
