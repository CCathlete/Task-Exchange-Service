package entities

// Represents a single task in the task management system.
type Task struct {
	ID          int     // The ID of the task.
	Description string  // Description of the task.
	AssignedTo  int     // The ID of the user the task is assigned to.
	Status      string  // "pending" or "completed".
	Price       float64 // Cost or reward for completing the task.
}
