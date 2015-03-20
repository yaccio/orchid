/*
Orchid is an application for running scripts locally and remotely for testing
and deployment. The application consists of a command line interface and a CI
server.
*/

package main

import (
	"flag"
	"fmt"
	"os"
)

/*
Entry point for the Orchid command line interface
*/
func main() {
	// Make sure some command is given
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Handle flags and arguments
	var path string
	flag.StringVar(&path, "-p", "orchid", "Specify the path to the config directory")
	flag.Parse()
	var args = flag.Args()

	actions := Actions{path}

	// Create logs dir if it does not exist
	os.Mkdir("orchid/logs", 0744)

	// Run job
	if args[0] == "run" {
		if len(args) != 2 {
			printUsage()
			return
		}

		jobId := args[1]
		actions.RunJob(jobId)
	}

	// Execute action
	if args[0] == "exec" {
		if len(args) != 2 {
			printUsage()
			return
		}

		actionId := args[1]
		actions.ExecuteAction(actionId)
	}

	// List
	if args[0] == "list" {
		if len(args) != 2 {
			printUsage()
			return
		}

		if args[1] == "jobs" {
			// List jobs
			actions.ListJobs()
		} else if args[1] == "actions" {
			// List actions
			actions.ListActions()
		} else if args[1] == "machines" {
			// List machines
			actions.ListMachines()
		} else if args[1] == "scripts" {
			// List scripts
			actions.ListScripts()
		} else if args[1] == "logs" {
			// List logs
			actions.ListLogs()
		} else {
			printUsage()
		}
	}

	// Get log output
	if args[0] == "logs" {
		if len(args) != 2 {
			printUsage()
			return
		}

		logId := args[1]
		actions.GetLogOutput(logId)
	}

	// SSH into a given machine
	if args[0] == "ssh" {
		if len(args) != 2 {
			printUsage()
			return
		}

		machineId := args[1]
		actions.SSH(machineId)
	}
}

/*
Prints a help message, explaining how to use the application
*/
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("- list jobs\t// List all configured jobs")
	fmt.Println("- list machines\t// List all configured machines")
	fmt.Println("- list scripts\t// List all configured scripts")
	fmt.Println("- list logs\t// List all stored logs")
	fmt.Println("- run <job id>\t// Run the job with the given id")
	fmt.Println("- exec <action id>\t// Execute the action with the given id")
	fmt.Println("- logs <log id>\t// Tail the log with the given id")
	fmt.Println("- ssh <machine id>\t// SSH into the machine with the given id")
}
