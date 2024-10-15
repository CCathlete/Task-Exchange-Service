package infrastructure

import (
	"aTES/core/entities"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// TODO: Add a section for the creation of the database if it doesn't exist.
func InitDB(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Creating a new task and returning its taskID.
// TODO: Fix the time to include date only.
func CreateTask(db *sql.DB, description string, assignedTo int) (int, error) {
	price := rand.Float64()*20 + 20 // Random price between 20 and 40.
	query := `INSERT INTO tasks (description, assigned_to, status, price, creation_time, completion_time)` +
		`VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, description, assignedTo, "pending", price, time.Now().Format("YYYY-MM-DD HH:MM"), "")
	if err != nil {
		return 0, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(taskID), nil
}

// Getting tasks that are assigned to a specific user.
func GetTasks(db *sql.DB, userID int) ([]entities.Task, error) {
	query := `SELECT id, description, assigned_to, status, price FROM tasks WHERE assigned_to = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		var task entities.Task
		if err := rows.Scan(&task.ID, &task.Description, &task.AssignedTo, &task.Status, &task.Price); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func CreateUser(db *sql.DB, name, email, role, joinedAt string) (int, error) {
	var userID int
	return userID, nil
}

func GetUser(db *sql.DB, name, email, role string) (entities.User, error)
