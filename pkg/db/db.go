package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"database/sql"

	_ "modernc.org/sqlite"
)

// код для создания таблиц БД
var schema = `CREATE TABLE scheduler(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT '',
            title VARCHAR(32) NOT NULL DEFAULT '',
            comment TEXT NOT NULL DEFAULT '',
            repeat VARCHAR(128) NOT NULL DEFAULT ''
        );
        CREATE INDEX date_index ON scheduler(date);`

// наличие бд если нет создаем
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("os.Stat: %w", err)
		}
		path, err := filepath.Abs(dbFile)
		if err != nil {
			return fmt.Errorf("filepath.Abs: %w", err)
		}
		_, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		DB, err := sql.Open("sqlite", path)
		if err != nil {
			return fmt.Errorf("sql.Open: %w", err)
		}
		defer DB.Close()
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("DB.Exec: %w", err)
		}
	}
	return nil
}
