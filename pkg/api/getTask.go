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
  // плохой нейминг db, а не task, так как он возвращает задачу а не db
	db, err := db.UpdateTask(id, DB)
	if err != nil {
		// у тебя могут быть разные причины ошибок не только что задача не найдена
		// надо обработать все возможные ошибки
		// Например:
		//if errors.Is(err, sql.ErrNoRows) {
		//	services.Er(w, err, http.StatusNotFound)
		//} else{
		//	services.Er(w, err, http.StatusInternalServerError)
		//}

		services.Er(w, err, http.StatusNotFound)
		return
	}

	services.WriteJson(w, db, http.StatusOK)
}
