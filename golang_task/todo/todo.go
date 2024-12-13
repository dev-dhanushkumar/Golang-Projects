package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

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

func (t *Todos) Update(id int, task string, cat string, done int) error {
	ls := *t

	index := t.getIndexByID(id)
	if index == -1 {
		return errors.New("invalid ID")
	}

	if len(task) != 0 {
		ls[index].Task = task
	}

	if len(cat) != 0 {
		ls[index].Category = cat
	}

	// The Error occurs because we are trying to assign a value of type time.Time
	// to a variable of type *time.Time . To fix this we can use the new() function to create a new pointer
	// to time. Time and assing it to CompletedAt. By using &completedAt, we create a new pointer to time . Time that to the CompletedAtvariable
	if done == 0 {
		ls[index].Done = false
		ls[index].CompletedAt = nil
	} else if done == 1 {
		ls[index].Done = true
		completedAt := time.Now()
		ls[index].CompletedAt = &completedAt // create a new pointer to time.
	}
	return nil
}

// Delete will delete requested tast from slice Todos
func (t *Todos) Delete(id int) error {
	ls := *t
	index := t.getIndexByID(id)
	if index == -1 {
		return errors.New("invalid ID")
	}

	*t = append(ls[:index], ls[index+1:]...)
	return nil
}

// Load will read .todos.json file and update date into Todos slice
func (t *Todos) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		return err
	}

	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	// Update nextId based on the loaded tasks
	if len(*t) > 0 {
		maxID := (*t)[0].ID
		for _, todo := range *t {
			if todo.ID > maxID {
				maxID = todo.ID
			}
		}
		nextID = maxID + 1
	}
	return nil
}

// Store will write Todos silce data into .todos.json file
func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Print will print out the current todo tasks
func (t *Todos) Print(status int, cat string) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	requestedTodos := []item{}

	for _, todo := range *t {
		if status == 1 {
			if todo.Done {
				requestedTodos = append(requestedTodos, todo)
			}
		}

		if status == 0 {
			if !todo.Done {
				requestedTodos = append(requestedTodos, todo)
			}
		}

		if status != 1 && status != 0 {
			requestedTodos = append(requestedTodos, todo)
		}
	}

	requestedCatTodos := []item{}

	for _, todo := range requestedTodos {
		if strings.ToLower(todo.Category) == strings.ToLower(cat) || cat == "" {
			requestedCatTodos = append(requestedCatTodos, todo)
		}
	}

	for _, item := range requestedCatTodos {
		task := item.Task
		done := "No"
		completedAt := ""

		if item.Done {
			task = fmt.Sprintf("%s", item.Task)
			done = "\u2705"
		}

		if item.CompletedAt != nil {
			completedAt = item.CreatedAt.Format("2006-01-02")
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", item.ID)},
			{Text: item.Category},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format("2006-01-02")},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: ""},
		{Align: simpletable.AlignLeft, Span: 5, Text: fmt.Sprintf("You have %d pending todos", t.CountPending())},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

// CountPending() will print out the pending tasks
func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}
	return total
}

// getIndexByID returns the index from a given item's id
func (t *Todos) getIndexByID(id int) int {
	index := -1
	for i, todo := range *t {
		if todo.ID == id {
			index = i
			break
		}
	}
	return index
}
