/*
A pipeline wraps the executable part of a job including command execution and
logging
*/

package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

/*
Type defining the pipeline
*/
type Pipeline struct {
	Cmds []*exec.Cmd
	Log  Log
	File *os.File
}

/*
Run/execute the pipeline, executing the commands it containes sequentially,
aborting if an error is encountered. This includes updating the logs file.
*/
func (p Pipeline) Run(path string) {
	// Always close the file after use
	defer p.File.Close()

	var err error

	// Write to the logs file that the job has started
	p.Log, err = p.Log.start(path)
	if err != nil {
		p.File.Close()
		p.Log.error(path, p.File)
		return
	}

	// Run the commands
	for _, cmd := range p.Cmds {
		err = cmd.Start()
		if err != nil {
			return
		}
		err = cmd.Wait()
		if err != nil {
			return
		}
	}

	// Write to the logs file that the job has finished, terminating
	// any tails following the log, once the job has finished
	p.Log, _ = p.Log.finish(path, p.File)
	//TODO find a way of handling the error that might be thrown
}

/*
Build a pipeline from a job
*/
func buildPipeline(path, jobId string, log Log) (Pipeline, error) {
	setup, err := loadSetup(path)
	if err != nil {
		return Pipeline{}, err
	}

	var job Job
	jobFound := false
	for _, j := range setup.Jobs {
		if j.Id == jobId {
			job = j
			jobFound = true
			break
		}
	}

	if !jobFound {
		return Pipeline{}, errors.New("Job not found")
	}

	logPath := fmt.Sprintf("%s/logs/%s", path, log.Id)
	outfile, err := os.Create(logPath)
	if err != nil {
		return Pipeline{}, err
	}

	var pipeline Pipeline
	pipeline.File = outfile
	pipeline.Log = log
	for _, executable := range job.Pipeline {
		cmd, execErr := buildExecutable(path, executable, setup.Machines, log, outfile)
		if execErr != nil {
			return Pipeline{}, execErr
		}
		pipeline.Cmds = append(pipeline.Cmds, cmd)
	}

	return pipeline, nil
}

/*
Build a command executable by the OS from an executable as defined in the job
configuration
*/
func buildExecutable(path string, executable Executable, machines []Machine, log Log, file *os.File) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	script := path + "/scripts/" + executable.Script
	if executable.Machine == "local" {
		cmd = exec.Command("/bin/bash", script)
		cmd.Stdout = file
	} else {
		var machine Machine
		for _, m := range machines {
			if m.Id == executable.Machine {
				machine = m
				break
			}
		}

		sshCommand := fmt.Sprintf(
			"ssh %s@%s -p %s -i %s 'bash -s' < %s",
			machine.User,
			machine.Address,
			machine.Port,
			path+"/keys/"+machine.PrivateKey,
			script,
		)
		cmd = exec.Command("/bin/bash", "-c", sshCommand)
	}

	return cmd, nil
}
