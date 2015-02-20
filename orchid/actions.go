/*
Implementation of the Api for executing commands locally
*/

package main

import (
	"fmt"
	"github.com/ActiveState/tail"
	"strings"
)

type Actions struct {
	path string
}

/*
List all jobs
*/
func (a *Actions) ListJobs() {
	setup, err := loadSetup(a.path)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	for _, job := range setup.Jobs {
		fmt.Println(job.Id)
		for _, ex := range job.Pipeline {
			fmt.Printf("\t%s -> %s\n", ex.Machine, ex.Script)
		}
	}
}

/*
List all machines
*/
func (a *Actions) ListMachines() {
	setup, err := loadSetup(a.path)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	for _, machine := range setup.Machines {
		fmt.Println(machine.Id)
		fmt.Printf("\t%s@%s:%s (%s)\n", machine.User, machine.Address, machine.Port, machine.PrivateKey)
	}
}

/*
List all scripts
*/
func (a *Actions) ListScripts() {
	setup, err := loadSetup(a.path)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	for _, script := range setup.Scripts {
		fmt.Println(script)
	}
}

/*
List all existing logs stored locally
*/
func (a *Actions) ListLogs() {
	logs, err := loadLogs(a.path)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	fmt.Printf("%-20s\t%-20s\t%-20s\t%-32s\t%-32s\n", "Id", "Job", "Status", "Start", "End")
	for _, log := range logs {
		fmt.Printf("%-20s\t%-20s\t%-20s\t%-32s\t%-32s\n", log.Id, log.JobId, log.Status, log.StartTime, log.EndTime)
	}
}

/*
Execute the job with the given id locally
*/
func (a *Actions) RunJob(jobId string) {
	log := newLog(jobId)

	pipeline, err := buildPipeline(a.path, jobId, log)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	go func() {
		pipeline.Run(a.path)
	}()

	fmt.Println(log.Id)

	// Tail the log, ensuring the program does not terminate
	a.GetLogOutput(log.Id)
}

/*
Get the output stored locally in the log with the given id
*/
func (a *Actions) GetLogOutput(logId string) {
	// If the log id given is not full, search for the first log that
	// matches the id prefix
	if len(logId) < 16 {
		match := ""
		logs, err := loadLogs(a.path)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}

		for _, log := range logs {
			if strings.HasPrefix(log.Id, logId) {
				match = log.Id
				break
			}
		}

		// If no match, inform the user
		if match == "" {
			fmt.Println("ERROR: Log not found")
			return
		}

		logId = match

	}
	t, err := tail.TailFile(a.path+"/logs/"+logId, tail.Config{Follow: true})
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	for line := range t.Lines {
		if line.Text == "-----Finished-----" || line.Text == "-----Error-----" {
			break
		}
		fmt.Println(line.Text)
	}
}
