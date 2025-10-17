package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"final/pkg/api/services"
	"final/pkg/db"
)

func PutTask(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var dto db.TaskDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	jsonString := string(body)

	str, _, err := services.Check(json.NewDecoder(strings.NewReader(jsonString)))
	if err != nil {
		services.Er(w, err, http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&dto)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	dto.Date = str

	err = db.TaskPut(&dto, DB)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	services.WriteJson(w, map[string]string{}, http.StatusCreated)
}
