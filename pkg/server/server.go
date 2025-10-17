package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"final/pkg/api"
	"final/pkg/db"

	_ "modernc.org/sqlite"
)

// порт по умочанию
var port = "7540"
var pathDB = "./pkg/db/scheduler.db"

func Init(router *http.ServeMux, db *sql.DB) {
	router.Handle("/", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		api.NextDayHandler(w, r, db)
	})
	router.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		api.TaskHandler(w, r, db)
	})
	router.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		api.TasksHandler(w, r, db)
	})
	router.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		api.DoneHandler(w, r, db)
	})
}

func StartServer(logger *log.Logger) *http.Server {

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		pathDB = path
	}

	err := db.Init(pathDB)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	Init(router, db)
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	return &server
}
