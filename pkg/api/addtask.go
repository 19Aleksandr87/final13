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

	/*_, err = DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES($1, $2, $3, $4)",
		content.Date, content.Title, content.Comment, content.Repeat)
	if err != nil {
		services.Er(w, errors.New(`{"error":"не удалось записать в бд"}`), http.StatusInternalServerError)
		return
	}*/

	var lastID int
	row := DB.QueryRow("INSERT INTO scheduler (date, title, comment, repeat) VALUES($1, $2, $3, $4) RETURNING id",
		content.Date, content.Title, content.Comment, content.Repeat)
	err = row.Scan(&lastID)
	if err != nil {
		services.Er(w, errors.New(`{"error":"не удалось записать в бд"}`), http.StatusInternalServerError)
		return
	}

	s := strconv.Itoa(lastID)
	services.WriteJson(w, map[string]string{"id": s}, http.StatusCreated)

}
