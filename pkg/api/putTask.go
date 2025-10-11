package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"final/pkg/api/services"
	"final/pkg/db"
)

func PutTask(w http.ResponseWriter, r *http.Request) {
	var dto db.TaskDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		services.Er(w, err)
		return
	}
	jsonString := string(body)

	str, _, err := services.Check(json.NewDecoder(strings.NewReader(jsonString)))
	if err != nil {
		services.Er(w, err)
		return
	}
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&dto)
	if err != nil {
		services.Er(w, err)
		return
	}
	dto.Date = str

	err = db.TaskPut(&dto)
	if err != nil {
		services.Er(w, err)
		return
	}

	w.Write([]byte(`{}`))
}
