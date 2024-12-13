package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dev-dhanushkumar/golang-projects/mytask/cmd"
	"github.com/dev-dhanushkumar/golang-projects/mytask/todo"
)

func main() {
	todos := &todo.Todos{}
	flag.Parse()

	// Check whick subcommand was invoked
	switch flag.Arg(0) {
	case "":
		cmd.Help()
	case "help":
		cmd.Help()
	case "init":
		cmd.Init()
	case "add":
		cmd.RemindInit(todos)
		cmd.AddTask(todos, os.Args[2:])
	case "list":
		cmd.RemindInit(todos)
		cmd.ListTasks(todos, os.Args[2:])
	case "delete":
		cmd.RemindInit(todos)
		cmd.DeleteTask(todos, os.Args[2:])
	case "update":
		cmd.RemindInit(todos)
		cmd.UpdateTask(todos, os.Args[2:])
	default:
		fmt.Println("Invalid command. Please use 'help' command to see available commands.")
		os.Exit(1)
	}
}
