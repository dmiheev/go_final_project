package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	ID    int64  `json:"id,omitempty,string"`
	Error string `json:"error,omitempty"`
}

func Init(ts TaskService) {
	http.HandleFunc("/api/nextdate", ts.nextDayHandler)
	http.HandleFunc("/api/tasks", ts.tasksHandler)
	http.HandleFunc("/api/task/done", ts.taskDoneHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ts.getTaskHandler(w, r)
		case http.MethodPost:
			ts.taskHandler(w, r)
		case http.MethodPut:
			ts.updateTaskHandler(w, r)
		case http.MethodDelete:
			ts.taskDeleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func responseError(w http.ResponseWriter, message string, statusCode int) {
	response := Response{Error: message}
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
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, `{"error":"failed to write JSON"}`, http.StatusInternalServerError)
		return
	}
}
