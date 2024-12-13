package todo

import "time"

// todo struct
type item struct {
	ID          int
	Task        string
	Category    string
	Done        bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// []item - slice
type Todos []item

// NextID will keep track of the next available ID for a next task
var nextID int

// Add will add a new task to slice Todos
func (t *Todos) Add(task string, cat string) {
	todo := item{
		ID:          nextID,
		Task:        task,
		Category:    cat,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: nil, // set to nil insted of time.Time{}
	}

	// Increment nextID for the next task
	nextID++

	// add a new task to ToDos slice
	*t = append(*t, todo)
}
