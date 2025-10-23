package main

import (
	"final/pkg/server"
	"log"
	"os"
)

var setting = server.Settings{
	Port:   "7540",
	PathDB: "./pkg/db/scheduler.db",
}

func main() {
	// передавай путь к файлу или через переменные окружения или аргументы командной строки
	flog, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // слишком высокий уровень доступа 644 или 600
	if err != nil {
		log.Fatal(err)
	}
	//defer flog.Close() defer решил удалить так как он открывается один раз, а в случае завершения программы или ошибки автоматически закроется

	setting.Logger = log.New(flog, "serv", log.LstdFlags|log.Lshortfile)

	HTTPServer := server.StartServer(&setting)
	log.Fatal(HTTPServer.ListenAndServe()) // после Fatal не отработает дефер
}
