package server

import (
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

type Server struct {
	Loger  *log.Logger
	Server *http.Client
}

func Init(router *http.ServeMux) {
	router.Handle("/", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/nextdate", api.NextDayHandler)
	router.HandleFunc("/api/task", api.TaskHandler)
	router.HandleFunc("/api/tasks", api.TasksHandler)
	router.HandleFunc("/api/task/done", api.DoneHandler)
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

	router := http.NewServeMux()
	Init(router)
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	return &server
}
