package api

import (
	"database/sql"
	"errors"
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"
)

func DelTask(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		services.Er(w, errors.New(`{"error":"Method not allowed"}`), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		services.Er(w, errors.New(`{"error":"Не указан идентификатор"}`), http.StatusBadRequest)
		return
	}
	// вот тут правильно для создания тоже самое надо
	err := db.TaskDel(id, DB)
	if err != nil {
		services.Er(w, err, http.StatusNotFound)
		return
	}
	services.WriteJson(w, map[string]string{}, http.StatusOK)
}
