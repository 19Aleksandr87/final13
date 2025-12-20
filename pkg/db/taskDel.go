package db

import (
	"database/sql"
	"strconv"
)

func TaskDel(id string, DB *sql.DB) error {

	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = DB.Exec("DELETE FROM scheduler WHERE id = $1", i)
	if err != nil {
		return err
	}

	return nil
}
