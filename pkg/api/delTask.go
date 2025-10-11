package api

import (
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"
)

func DelTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	err := db.TaskDel(id)
	if err != nil {
		services.Er(w, err)
		return
	}
	w.Write([]byte(`{}`))
}
