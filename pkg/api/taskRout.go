package api

import (
	"net/http"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		AddTask(w, r)
		return
	case http.MethodGet:
		GetTask(w, r)
		return
	case http.MethodPut:
		PutTask(w, r)
		return
	case http.MethodDelete:
		DelTask(w, r)
		return
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}
