// Package cmd provides commanf-line sub-commands for the task application
package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dev-dhanushkumar/golang-projects/mytask/todo"
)

func UpdateTask(todos *todo.Todos, args []string) {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateId := updateCmd.Int("id", 0, "The id of todo to be updated.")
	updateCat := updateCmd.String("cat", "", "The to-be-updated category of todo")
	updateTask := updateCmd.String("task", "", "To to-be-updated content of todo")
	updateDone := updateCmd.Int("done", 2, "The to-be-updated status of todo")

	// Parse the argument for the "update" subcomand
	updateCmd.Parse(args)

	if *updateId == 0 {
		fmt.Println("Error: the --id flag is required for the 'update' subcommand.")
		os.Exit(1)
	}
	err := todos.Update(*updateId, *updateTask, *updateCat, *updateDone)
	if err != nil {
		log.Fatal(err)
	}

	err = todos.Store(GetJsonFile())
	if err != nil {
		log.Fatal(err)
	}

	todos.Print(2, "")
	fmt.Println("Todo item updated Successfully.")
}
