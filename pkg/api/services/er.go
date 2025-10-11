package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Er(w http.ResponseWriter, err error) {
	errorResponse := map[string]string{"error": err.Error()}
	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
		return
	}

	http.Error(w, string(jsonData), http.StatusBadRequest)
}
