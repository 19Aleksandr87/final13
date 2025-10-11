package db

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"
)

func TaskDel(id string) error {
	DB, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	if err != nil {
		return err
	}
	defer DB.Close()
	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = DB.Exec("DELETE FROM scheduler WHERE id = ?", i)
	if err != nil {
		return err
	}

	return nil
}
