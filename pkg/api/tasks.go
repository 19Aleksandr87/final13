package api

import (
	"database/sql"
	"net/http"

	"final/pkg/api/services"
	"final/pkg/db"

	"time"
)

type TasksResp struct {
	Tasks []*db.TaskDTO `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
	search := r.URL.Query().Get("search")
	if search != "" {
		t, err := time.Parse("02.01.2006", search)
		if err == nil {
			s := t.Format(services.FormatDate)
			tasks, err := db.DateTasks(s, 50, DB)
			if err != nil {
				services.Er(w, err, http.StatusInternalServerError)
				return
			}
			services.WriteJson(w, TasksResp{
				Tasks: tasks,
			}, http.StatusInternalServerError)
			return
		}
		err = nil
		tasks, err := db.SearchTasks(search, 50, DB)
		if err != nil {
			services.Er(w, err, http.StatusInternalServerError)
			return
		}
		services.WriteJson(w, TasksResp{
			Tasks: tasks,
		}, http.StatusInternalServerError)
		return
	}

	tasks, err := db.Tasks(50, DB) // в параметре максимальное количество записей
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		services.WriteJson(w, map[string][]any{"tasks": {}}, http.StatusOK)
		return
	}

	services.WriteJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusInternalServerError)
}
