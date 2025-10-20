package api

import (
	"database/sql"
	"errors"
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"
)

func GetTask(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		services.Er(w, errors.New(`{"error": "Не указан идентификатор"}`), http.StatusBadRequest)
		return
	}

	db, err := db.UpdateTask(id, DB)
	if err != nil {
		services.Er(w, err, http.StatusNotFound)
		return
	}

	services.WriteJson(w, db, http.StatusOK)

}
