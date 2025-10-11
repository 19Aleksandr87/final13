package api

import (
	"database/sql"
	"encoding/json"
	"final/pkg/api/services"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func AddTask(w http.ResponseWriter, r *http.Request) {

	DB, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Failed to open result file"}`, http.StatusInternalServerError)
		return
	}
	defer DB.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		services.Er(w, err)
		return
	}
	jsonString := string(body)

	str, content, err := services.Check(json.NewDecoder(strings.NewReader(jsonString)))
	if err != nil {
		services.Er(w, err)
		return
	}

	content.Date = str

	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)",
		sql.Named("date", content.Date), sql.Named("title", content.Title), sql.Named("comment", content.Comment), sql.Named("repeat", content.Repeat),
	)
	if err != nil {
		http.Error(w, `{"error":"не удалось записать в бд"}`, http.StatusInternalServerError)
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	s := strconv.Itoa(int(i))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": s})
}
