package infrastructure

import (
	"aTES/core/entities"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

/*
TODO: Make sure that the task table has task_id for each row and not just id.
Same with users table - user_id and accounting_records tables - record_id.
*/
func InitDB(config Config) (*sql.DB, error) {
	connectStringNoDB := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBSSLMode,
	)

	db, err := sql.Open("postgres", connectStringNoDB)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres server: %w", err)
	}
	defer db.Close()

	// Checking if our target DB exists. We're extracting a boolean from the query and an error
	// means there was a problem with the check.
	var itExists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s'", config.DBName)
	err = db.QueryRow(query).Scan(&itExists)
	if err != nil {
		return nil, fmt.Errorf("error checking if database exists: %w", err)
	}

	// We create the DB if it doesn't exist.
	if !itExists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
		if err != nil {
			return nil, fmt.Errorf("couldn't create the database: %w", err)
		}
		log.Printf("Database %s created.\n", config.DBName)
	} else {
		log.Printf("Database %s already exists.\n", config.DBName)
	}

	// Connecting to the target db. We'll use the string created earlier and att dbname to it.
	connectStringWithDB := fmt.Sprintf("%s dbname=%s", connectStringNoDB, config.DBName)

	db, err = sql.Open("postgres", connectStringWithDB)
	if err != nil {
		return nil, fmt.Errorf("error connecting to target database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection was successful but DB is not responding: %w", err)
	}

	return db, nil
}

// Creating a new task and returning its taskID.
func CreateTask(db *sql.DB, description string, assignedTo int) (int, error) {
	price := rand.Float64()*20 + 20 // Random price between 20 and 40.
	query := `INSERT INTO tasks (description, assigned_to, status, price, creation_time, completion_time)` +
		`VALUES ($1, $2, $3, $4, $5, $6)`

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

func UpdateTask(db *sql.DB, description, status string, taskID, assignedTo int, price float64, isCompleted bool) error {
	query := `
		UPDATE tasks
		SET description = $1,
			status = $2,
			assigned_to = $3,
			price = $4,
			is_completed = $5
		WHERE task_id = $6
		`
	_, err := db.Exec(query, description, status, assignedTo, price, isCompleted, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Getting tasks that are assigned to a specific user.
func GetTasks(db *sql.DB, userID int) ([]entities.Task, error) {
	query := `SELECT task_id, description, assigned_to, status, price FROM tasks WHERE assigned_to = $1`
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
	query := `
	INSERT INTO users (name, email, role, joined_at)
	VALUES ($1, $2, $3, $4)
	RETURNING user_id
	`
	err := db.QueryRow(query, name, email, role, joinedAt).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("there was an issue with creating the user: %w", err)
	}

	return userID, nil
}

// We use ID or cross role, name and email to find the correct user.
func GetUser(db *sql.DB, userIDp *int, name, email, role string) (entities.User, error) {
	var user entities.User
	var query string
	var err error

	if userIDp != nil {
		query = `
		SELECT user_id, name, email, role, joined_at 
		FROM users 
		WHERE id = $1
		`
		err = db.QueryRow(query, *userIDp).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.JoinedAt)
	} else {
		// If userID is not provided we cross other fields.
		query = `
		SELECT user_id, name, email, role, joined_at 
		FROM users 
		WHERE (name = $1 AND email = $2)
		OR (email = $2 AND role = $3)
		`
		err = db.QueryRow(query, name, email, role).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.JoinedAt)
	}

	if err != nil {
		return entities.User{}, fmt.Errorf("issue with getting user info from the DB: %w", err)
	}

	return user, nil
}

func UpdateUser(db *sql.DB, userID int, name, email, role string) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, role = $3
	WHERE user_id = $4
	`

	_, err := db.Exec(query, name, email, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func CreateAccountingRecord(db *sql.DB, userID, taskID int, status string, assignedTo int, price float64, isCompleted bool) (int, error) {
	var recordID int

	query := `
	INSERT INTO accounting_records (user_id, task_id, status, assigned_to, price, is_completed)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING record_id
	`
	err := db.QueryRow(query, userID, taskID, status, assignedTo, price, isCompleted).Scan(&recordID)
	if err != nil {
		return 0, fmt.Errorf("failed to create an accounting record: %w", err)
	}

	return recordID, nil
}
