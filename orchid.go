/*
Orchid is an application for running scripts locally and remotely for testing
and deployment. The application consists of a command line interface and a CI
server.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/yaccio/orchid/ci_server"
	"github.com/yaccio/orchid/core"
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
	var remote bool
	flag.BoolVar(&remote, "r", false, "Run the command on the remote CI server")
	var path string
	flag.StringVar(&path, "-p", "ci", "Give a specified path to the config directory")
	flag.Parse()
	var args = flag.Args()

	// Run CI server (this blocks)
	if args[0] == "ci" {
		if len(args) != 1 {
			printUsage()
			return
		}
		ci_server.Start()
		return

	}

	// Get local or remote api depending on flags
	var api core.Api
	var err error
	if remote {
		api, err = core.NewRemoteApi(path)
	} else {
		api, err = core.NewLocalApi(path)
	}
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	// Run job
	if args[0] == "run" {
		if len(args) != 2 {
			printUsage()
			return
		}

		jobId := args[1]

		var log core.Log
		log, err = api.RunJob(jobId)
		fmt.Println(log.Id)

		// Tail the log, ensuring the program does not terminate
		err = api.GetLogOutput(log.Id, os.Stdout)
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

		logs, err := api.ListLogs()
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

		err = api.GetLogOutput(logId, os.Stdout)
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
	fmt.Println("- Run CI server:\torchid ci")
	fmt.Println("- Run job:\torchid [-r] run <job id>")
	fmt.Println("- List logs:\torchid [-r] logs")
	fmt.Println("- Get log:\torchid [-r] logs <log id>")
	fmt.Println("The -r flag runs the command on the remote CI server")
}
