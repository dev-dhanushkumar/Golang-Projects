package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/dev-dhanushkumar/golang-projects/mytask/todo"
)

func DeleteTask(todos *todo.Todos, args []string) {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// If no --id=1 flag defined todo will default to 0
	// but if --id is present but didn't set any value an error will be thrown
	deleteID := deleteCmd.Int("id", 0, "The id of todo to be deleted")

	// Parse the argument for the "delete" subcommand
	deleteCmd.Parse(args)

	err := todos.Delete(*deleteID)
	if err != nil {
		log.Fatal(err)
	}

	err = todos.Store(GetJsonFile())
	if err != nil {
		log.Fatal(err)
	}

	todos.Print(2, "")
	fmt.Println("Todo item deleted successfully.")
}
