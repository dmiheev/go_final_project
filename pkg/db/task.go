package db

import (
	"fmt"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    int64  `json:"id,omitempty,string"`
	Error string `json:"error,omitempty"`
}

type TasksResp struct {
	Tasks []Task `json:"tasks"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	dbConn, _ = open()
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := dbConn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId: %w", err)
	}
	task.ID = id

	return id, nil
}

func GetTasks(search string) ([]Task, error) {
	var query string
	var args []interface{}

	if search != "" {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE date = ?
				ORDER BY date
				LIMIT 50
			`
			args = append(args, parsedDate.Format("20060102"))
		} else {
			searchPattern := "%" + search + "%"
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE title LIKE ? OR comment LIKE ?
				ORDER BY date
				LIMIT 50
			`
			args = append(args, searchPattern, searchPattern)
		}
	} else {
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date
			LIMIT 50
		`
	}

	dbConn, _ = open()
	rows, err := dbConn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to parse tasks: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

func GetTask(id int64) (*Task, error) {
	db, _ := open()
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`
	row := db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	db, _ := open()
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	db, _ := open()
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateTaskDate(id int64, nextDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	db, _ := open()
	res, err := db.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
