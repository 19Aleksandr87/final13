package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Er(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}
	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	http.Error(w, string(jsonData), code)
}
