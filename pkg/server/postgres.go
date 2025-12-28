package server

import (
	"database/sql"
	"final/pkg/db"
)

type postgresType struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func (p *postgresType) init(DB *sql.DB) error {
	name := "scheduler"
	datname := ""
	err := DB.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = $1", name).Scan(&datname)
	if err != nil {
		err = db.Init(DB)
		return err
	}
	return nil
}

var postgres = postgresType{
	//host:     "localhost",
	host:     "db",
	port:     5432,
	user:     "postgres",
	password: "postgres",
	dbname:   "postgres",
}
