package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"final/pkg/api/services"
	"final/pkg/db"
)

func DoneHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	if r.Method != http.MethodPost {
		services.Er(w, errors.New(`{"error":"Method not allowed"}`), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		services.Er(w, errors.New(`{"error": "Не указан идентификатор"}`), http.StatusNotFound)
		return
	}
	dto, err := db.UpdateTask(id, DB)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}

	if dto.Repeat != "" {
		now := time.Now()
		dstart, err := time.Parse(services.FormatDate, dto.Date)
		if err != nil {
			services.Er(w, err, http.StatusInternalServerError)
		}
		if now.Before(dstart) {
			now = dstart
		}
		tS := now.AddDate(0, 0, 1).Format(services.FormatDate)
		str, err := services.NextDate(tS, dto.Date, dto.Repeat)
		if err != nil {
			services.Er(w, err, http.StatusInternalServerError)
			return
		}
		dto.Date = str
		err = db.TaskPut(dto, DB)
		if err != nil {
			services.Er(w, err, http.StatusInternalServerError)
			return
		}
		services.WriteJson(w, map[string]string{}, http.StatusCreated)
		return
	}
	err = db.TaskDel(dto.ID, DB)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	services.WriteJson(w, map[string]string{}, http.StatusCreated)
}
