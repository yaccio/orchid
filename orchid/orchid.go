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
        "path/filepath"
        "log"
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

        currentdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
                log.Fatal("Can't determine full path!")
        }
        path = currentdir + "/" + path

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

	// Copy files/directories from one machine to another
	if args[0] == "scp" {
		if len(args) != 3 {
			printUsage()
			return
		}

		from := args[1]
		to := args[2]
		actions.SCP(from, to)
	}

	// Mount a remote directory locally
	if args[0] == "mount" {
		if len(args) != 4 {
			printUsage()
			return
		}

                err := actions.Mount(args[1],args[2],args[3])
                if err != nil {
                        fmt.Printf("Ouch, got error %#v, is the directory already mounted?",err)
                }
	}
	// Unmount a remote directory locally
	if args[0] == "unmount" {
		if len(args) != 2 {
			printUsage()
			return
		}

		err := actions.Unmount(args[1])
                if err != nil {
                        fmt.Printf("Ouch, got error %#v, is the directory mounted?",err)
                }
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
	fmt.Println("- scp <machine id>:<path> <machine id>:<path>\t// Copy files/directories from one machine to another. Only one of the machines can be specified. The other must be a path to a local file / directory without ':'")
        fmt.Println("- mount <machine id> <remote path> <local path>\t// Mount a remote directory (to which you have read access) locally")
        fmt.Println("- unmount <local path>\t// Unmount a previously Mount'ed directory")
}
