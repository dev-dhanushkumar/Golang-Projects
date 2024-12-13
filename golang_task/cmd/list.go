// Package cmd provides command-line subcommands for the task application
package cmd

import (
	"flag"

	"github.com/dev-dhanushkumar/golang-projects/mytask/todo"
)

func ListTasks(todos *todo.Todos, args []string) {
	// Define the list subcommand to list todo items
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listDone := listCmd.Int("done", 2, "The staus of todo to be printed")
	listCat := listCmd.String("cat", "", "The category of tasks to be listed")

	// Parse the argument for the "list" subcommand
	listCmd.Parse(args)
	todos.Print(*listDone, *listCat)
}
