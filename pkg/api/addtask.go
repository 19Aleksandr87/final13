package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final/pkg/api/services"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func AddTask(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	jsonString := string(body)

	str, content, err := services.Check(json.NewDecoder(strings.NewReader(jsonString)))
	if err != nil {
		services.Er(w, err, http.StatusBadRequest)
		return
	}

	content.Date = str

	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)",
		sql.Named("date", content.Date), sql.Named("title", content.Title), sql.Named("comment", content.Comment), sql.Named("repeat", content.Repeat),
	)
	if err != nil {
		services.Er(w, errors.New(`{"error":"не удалось записать в бд"}`), http.StatusInternalServerError)
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		services.Er(w, err, http.StatusInternalServerError)
		return
	}
	s := strconv.Itoa(int(i))
	services.WriteJson(w, map[string]string{"id": s}, http.StatusCreated)

}
