package api

import (
	"database/sql"
	"net/http"
)

func TaskHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		AddTask(w, r, DB)
	case http.MethodGet:
		GetTask(w, r, DB)
	case http.MethodPut:
		PutTask(w, r, DB)
	case http.MethodDelete:
		DelTask(w, r, DB)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
