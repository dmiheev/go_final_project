package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

const (
	DateFormat = "20060102"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPost:
			taskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		case http.MethodDelete:
			taskDeleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func responseError(w http.ResponseWriter, message string, statusCode int) {
	response := db.Response{Error: message}
	writeJSON(w, response, statusCode)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"failed to serialize JSON"}`, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
