package api

import (
	"net/http"
	"time"

	"final/pkg/api/services"
	"final/pkg/db"
)

func DoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	dto, err := db.UpdateTask(id)
	if err != nil {
		services.Er(w, err)
		return
	}

	if dto.Repeat != "" {
		tS := time.Now().AddDate(0, 0, 1).Format(services.FormatDate)
		str, err := services.NextDate(tS, dto.Date, dto.Repeat)
		if err != nil {
			services.Er(w, err)
			return
		}
		dto.Date = str
		err = db.TaskPut(dto)
		if err != nil {
			services.Er(w, err)
			return
		}
		w.Write([]byte(`{}`))
		return
	}
	err = db.TaskDel(dto.ID)
	if err != nil {
		services.Er(w, err)
		return
	}
	w.Write([]byte(`{}`))
}
