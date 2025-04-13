package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/utils"
)

type TaskService struct {
	store db.Storage
}

func NewTaskService(store db.Storage) TaskService {
	return TaskService{store: store}
}

func (t TaskService) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50"
	}

	tasks, err := t.store.GetTasks(search, limit)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := db.TasksResp{Tasks: tasks}
	writeJSON(w, response, http.StatusOK)
}

func (t TaskService) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		responseError(w, "failed to read the request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseError(w, "failed to deserialize JSON", http.StatusBadRequest)
		return
	}

	if err = task.Validate(); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.store.UpdateTask(&task)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func (t TaskService) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		responseError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		responseError(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTask(parsedId)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := db.Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	writeJSON(w, response, http.StatusOK)
}

func (t TaskService) taskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		responseError(w, "failed to read the request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseError(w, "failed to deserialize JSON", http.StatusBadRequest)
		return
	}

	if err = task.Validate(); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = t.store.AddTask(&task); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{ID: task.ID}
	writeJSON(w, response, http.StatusOK)
}

func (t TaskService) taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		responseError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		responseError(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTask(parsedId)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = t.store.DeleteTask(parsedId)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{}, http.StatusOK)
		return
	}

	parsedDate, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		responseError(w, "invalid date format, expected YYYYMMDD", http.StatusBadRequest)
	}

	nextDate, err := utils.NextDate(parsedDate, task.Date, task.Repeat)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.Date = nextDate

	err = t.store.UpdateTask(task)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func (t TaskService) taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		responseError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		responseError(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	err = t.store.DeleteTask(parsedId)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func (t TaskService) nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "missing required parameters: now, date or repeat", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(utils.DateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid format for 'now': "+nowStr, http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
