package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func ValidateTask(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Title == "" {
		return errors.New("task title is required")
	}

	if task.Date == "" {
		task.Date = today.Format(DateFormat)
		return nil
	}

	parsedDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	if parsedDate.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format(DateFormat)
			return nil
		}

		nextDate, err := NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = nextDate
		return nil
	}

	if task.Repeat != "" {
		_, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.GetTasks(search)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := db.TasksResp{Tasks: tasks}
	writeJSON(w, response, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = ValidateTask(&task); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	task, err := db.GetTask(parsedId)
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

func taskHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = ValidateTask(&task); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = db.AddTask(&task); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := db.Response{ID: task.ID}
	writeJSON(w, response, http.StatusOK)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
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

	task, err := db.GetTask(parsedId)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(parsedId)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{}, http.StatusOK)
		return
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		responseError(w, "invalid date format, expected YYYYMMDD", http.StatusBadRequest)
	}

	nextDate, err := NextDate(parsedDate, task.Date, task.Repeat)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTaskDate(parsedId, nextDate)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteTask(parsedId)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
