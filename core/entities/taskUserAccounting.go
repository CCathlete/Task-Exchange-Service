package entities

// Represents a single task in the task management system.
type Task struct {
	ID             int     // The ID of the task.
	Description    string  // Description of the task.
	AssignedTo     int     // The ID of the user the task is assigned to.
	Status         string  // Pending/completed/cancelled/started.
	Price          float64 // Cost or reward for completing the task.
	CreationTime   string  // Timestamp of creation time.
	CompletionTime string  // Timestamp of completion time.
	LastUpdated    string  // Timestamp of last update time.
}

type User struct {
	ID          int
	Name        string
	Email       string
	Role        string
	Balance     float64
	JoinedAt    string // Date of joining the company.
	LeftAt      string // Date of departure, empty list if currently employed.
	LastUpdated string // Timestamp of last update time.
}

type AccountingRecord struct {
	ID           int     // A unique ID for this record.
	UserID       int     // The ID of the user associated with this record.
	TaskID       int     // The ID of the task associated with this reduction/ payment.
	Amount       float64 // Negative for reduction and positive for payment.
	Status       string  // Assigned/ Completed.
	CreationTime string  // Timestamp of the creation time of this record.
	LastUpdated  string  // Timestamp of last update time.
}
