package entities

// Represents a single task in the task management system.
type Task struct {
	TaskID         int     `gorm:"primaryKey;autoIncrement"` // The ID of the task.
	Description    string  `gorm:"type:text"`                // Description of the task.
	AssignedTo     int     `gorm:"index;foreignKey:UserID"`  // The ID of the user the task is assigned to.
	Status         string  `gorm:"type:varchar(50)"`         // Pending/completed/cancelled/started.
	Price          float64 `gorm:"type:decimal(10, 2)"`      // Cost or reward for completing the task.
	CreationTime   string  `gorm:"type:timestamp"`           // Timestamp of creation time.
	CompletionTime string  `gorm:"type:timestamp"`           // Timestamp of completion time.
	LastUpdated    string  `gorm:"type:timestamp"`           // Timestamp of last update time.
}

type User struct { // We send this in http requests for the authorisation system to store.
	UserID      int     `json:"user_id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Role        string  `json:"role"`
	Balance     float64 `json:"balance"`
	JoinedAt    string  `json:"joined_at"`    // Date of joining the company.
	LeftAt      string  `json:"left_at"`      // Date of departure, empty list if currently employed.
	LastUpdated string  `json:"last_updated"` // Timestamp of last update time.
}

type TaskFinanceInfo struct {
	TaskID       int     `gorm:"primaryKey;autoIncrement"` // The ID of the task associated with this reduction/ payment.
	UserID       int     `gorm:"index;foreignKey:UserID"`  // The ID of the user associated with this record.
	kmount       float64 `gorm:"type:decimal(10, 2)"`      // Negative for reduction and positive for payment.
	Status       string  `gorm:"type:varchar(50)"`         // Assigned/ Completed.
	CreationTime string  `gorm:"type:timestamp"`           // Timestamp of the creation time of this record.
	LastUpdated  string  `gorm:"type:timestamp"`           // Timestamp of last update time.
}

// TODO: Change the type of the money amounts into something from the decimal library.
// TODO: An accounting record should be a monthly balance, amount of tasks assigned, amount of task started and amount of tasks completed.
