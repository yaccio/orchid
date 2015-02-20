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
	flag.StringVar(&path, "-p", "ci", "Give a specified path to the config directory")
	flag.Parse()
	var args = flag.Args()

	actions := Actions{path}

	// Run job
	if args[0] == "run" {
		if len(args) != 2 {
			printUsage()
			return
		}

		jobId := args[1]

		var log Log
		log, err := actions.RunJob(jobId)
		fmt.Println(log.Id)

		// Tail the log, ensuring the program does not terminate
		err = actions.GetLogOutput(log.Id)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	}

	// List jobs
	if args[0] == "list" {
		if len(args) != 1 {
			printUsage()
			return
		}

		logs, err := actions.ListLogs()
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Printf("%-20s\t%-20s\t%-20s\t%-32s\t%-32s\n", "Id", "Job", "Status", "Start", "End")
		for _, log := range logs {
			fmt.Printf("%-20s\t%-20s\t%-20s\t%-32s\t%-32s\n", log.Id, log.JobId, log.Status, log.StartTime, log.EndTime)
		}
	}

	// Get log
	if args[0] == "logs" {
		if len(args) != 2 {
			printUsage()
			return
		}

		logId := args[1]

		err := actions.GetLogOutput(logId)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	}
}

/*
Prints a help message, explaining how to use the application
*/
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("- Run job:\torchid run <job id>")
	fmt.Println("- List logs:\torchid logs")
	fmt.Println("- Get log:\torchid logs <log id>")
}
