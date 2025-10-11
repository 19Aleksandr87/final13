package api

import (
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	db, err := db.UpdateTask(id)
	if err != nil {
		services.Er(w, err)
		return
	}

	services.WriteJson(w, db)

}
