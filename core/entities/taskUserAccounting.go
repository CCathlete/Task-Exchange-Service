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

type User struct {
	UserID      int     `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"type:varchar(100)"`
	Email       string  `gorm:"type:varchar(100);uniqueIndex"`
	Role        string  `gorm:"type:varchar(50)"`
	Balance     float64 `gorm:"type:decimal(10, 2)"`
	JoinedAt    string  `gorm:"type:date"`      // Date of joining the company.
	LeftAt      string  `gorm:"type:date"`      // Date of departure, empty list if currently employed.
	LastUpdated string  `gorm:"type:timestamp"` // Timestamp of last update time.
}

type AccountingRecord struct {
	RecordID     int     `gorm:"primaryKey;autoIncrement"` // A unique ID for this record.
	UserID       int     `gorm:"index;foreignKey:UserID"`  // The ID of the user associated with this record.
	TaskID       int     `gorm:"index;foreignKey:TaskID"`  // The ID of the task associated with this reduction/ payment.
	Amount       float64 `gorm:"type:decimal(10, 2)"`      // Negative for reduction and positive for payment.
	Status       string  `gorm:"type:varchar(50)"`         // Assigned/ Completed.
	CreationTime string  `gorm:"type:timestamp"`           // Timestamp of the creation time of this record.
	LastUpdated  string  `gorm:"type:timestamp"`           // Timestamp of last update time.
}
