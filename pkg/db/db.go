package db

import (
	"fmt"

	"database/sql"
)

// код для создания таблиц БД
var schema = `CREATE TABLE scheduler(
            id SERIAL PRIMARY KEY,
            date CHAR(8) NOT NULL DEFAULT '',
            title VARCHAR(32) NOT NULL DEFAULT '',
            comment TEXT NOT NULL DEFAULT '',
            repeat VARCHAR(128) NOT NULL DEFAULT ''
        );
        CREATE INDEX date_index ON scheduler(date);`

// наличие бд если нет создаем
func Init(DB *sql.DB) error {
	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("DB.Exec: %w", err)
	}
	return nil
}
