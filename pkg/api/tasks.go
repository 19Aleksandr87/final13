package api

import (
	"encoding/json"
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.TaskDTO `json:"tasks"`
}

// func tasksHandler(w http.ResponseWriter) {
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		services.Er(w, err)
		return
	}
	if tasks == nil {
		w.Write([]byte(`{"tasks":[]}`))
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsData)
}
