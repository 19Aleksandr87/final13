package api

import (
	"net/http"

	"final/pkg/api/services"
)

type nextDate struct {
	Now    string `json:"now"`
	Date   string `json:"date"`
	Repeat string `json:"repeat"`
}

var nd nextDate

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nd.Now = r.FormValue("now")
	nd.Date = r.FormValue("date")
	nd.Repeat = r.FormValue("repeat")
	str, err := services.NextDate(nd.Now, nd.Date, nd.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(str))
}
