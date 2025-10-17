package main

import (
	"final/pkg/server"
	"log"
	"os"
)

func main() {
	flog, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()

	logger := log.New(flog, "serv", log.LstdFlags|log.Lshortfile)
	HTTPServer := server.StartServer(logger)
	log.Fatal(HTTPServer.ListenAndServe())
}
