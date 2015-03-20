/*
Definition of and methods for loading and validating the configuration files
concerned with jobs, machines, and scripts
*/

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
Type defining the complete setup of jobs, machines, and scripts
*/
type Setup struct {
	Machines []Machine
	Jobs     []Job
	Actions  []Action
	Scripts  []string
}

/*
Type defining a machine configuration
*/
type Machine struct {
	Id         string
	Address    string
	Port       string
	User       string
	PrivateKey string
}

/*
Type defining a job configuration
*/
type Job struct {
	Id       string
	Pipeline []Executable
}

/*
Type defining an executable (part of a job)
*/
type Executable struct {
	Machine string
	Script  string
	Args    []string
}

/*
Type defining an action
*/
type Action struct {
	Id      string
	Machine string
	Command string
}

/*
Load the configuration files concerned with the setup
*/
func loadSetup(path string) (Setup, error) {
	machines, machineErr := loadMachines(path)
	if machineErr != nil {
		return Setup{}, machineErr
	}

	jobs, jobErr := loadJobs(path)
	if jobErr != nil {
		return Setup{}, jobErr
	}

	actions, actionErr := loadActions(path)
	if actionErr != nil {
		return Setup{}, actionErr
	}

	scripts, scriptErr := loadDir(path + "/scripts")
	if scriptErr != nil {
		return Setup{}, scriptErr
	}

	keys, keyErr := loadDir(path + "/keys")
	if keyErr != nil {
		return Setup{}, keyErr
	}

	machineValidationErr := validateMachines(machines, keys, path)
	if machineValidationErr != nil {
		return Setup{}, machineValidationErr
	}

	jobValidationErr := validateJobs(jobs, machines, scripts, path)
	if jobValidationErr != nil {
		return Setup{}, jobValidationErr
	}

	actionValidationErr := validateActions(actions, machines)
	if actionValidationErr != nil {
		return Setup{}, actionValidationErr
	}

	setup := Setup{
		Machines: machines,
		Jobs:     jobs,
		Actions:  actions,
		Scripts:  scripts,
	}
	return setup, nil
}

/*
Load the configuration files concerned with machines
*/
func loadMachines(path string) ([]Machine, error) {
	machines := &[]Machine{}
	data, err := ioutil.ReadFile(path + "/machines.json")
	if err != nil {
		return []Machine{}, err
	}

	err = json.Unmarshal(data, &machines)
	if err != nil {
		return []Machine{}, err
	}

	return *machines, nil
}

/*
Load the configuration files concerned with jobs
*/
func loadJobs(path string) ([]Job, error) {
	jobs := &[]Job{}
	data, err := ioutil.ReadFile(path + "/jobs.json")
	if err != nil {
		return []Job{}, err
	}

	err = json.Unmarshal(data, &jobs)
	if err != nil {
		return []Job{}, err
	}

	return *jobs, nil
}

/*
Load the configuration files concerned with actions
*/
func loadActions(path string) ([]Action, error) {
	actions := &[]Action{}
	data, err := ioutil.ReadFile(path + "/actions.json")
	if err != nil {
		return []Action{}, err
	}

	err = json.Unmarshal(data, &actions)
	if err != nil {
		return []Action{}, err
	}

	return *actions, nil
}

/*
Helper method for loading the names of all files in a single directory.
Used for loading scripts and keys
*/
func loadDir(path string) ([]string, error) {
	var files = []string{}
	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			// Not allowed to visit the file. Continue walking
			return nil
		}
		if fi.IsDir() {
			// Ignore directories
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, nil
}

/*
Validate the machine configuration
*/
func validateMachines(machines []Machine, keys []string, path string) error {

	for _, machine := range machines {
		if machine.Id == "" {
			return errors.New("Machine config invalid: Each machine must have a non-empty id")
		}
		if machine.Address == "" {
			return errors.New("Machine config invalid: Machine '" + machine.Id + "' must have a non-empty Address")
		}
		if machine.Port == "" {
			return errors.New("Machine config invalid: Machine '" + machine.Id + "' must have a non-empty Port")
		}
		if machine.User == "" {
			return errors.New("Machine config invalid: Machine '" + machine.Id + "' must have a non-empty User")
		}
		if machine.PrivateKey == "" {
			return errors.New("Machine config invalid: Machine '" + machine.Id + "' must have a non-empty PrivateKey")
		}

		pathLength := len(path + "/keys")
		found := false
		for _, key := range keys {
			if machine.PrivateKey == key[pathLength+1:] {
				found = true
				break
			}
		}
		if !found {
			return errors.New("Machine config invalid: Machine '" + machine.Id + "' contains reference to unknown PrivateKey")
		}
	}

	return nil
}

/*
Validate the job configuration
*/
func validateJobs(jobs []Job, machines []Machine, scripts []string, path string) error {
	for _, job := range jobs {
		if job.Id == "" {
			return errors.New("Job config invalid: Each job must have a non-empty id")
		}
		if len(job.Pipeline) == 0 {
			return errors.New("Job config invalid: Job '" + job.Id + "' must have a non-empty Pipeline")
		}

		for _, executable := range job.Pipeline {
			machineFound := false
			for _, machine := range machines {
				if executable.Machine == machine.Id || executable.Machine == "local" {
					machineFound = true
					break
				}
			}
			if !machineFound {
				return errors.New("Job config invalid: Job '" + job.Id + "' contains a reference to one or more unknown machines")
			}

			pathLength := len(path + "/scripts")
			scriptFound := false
			for _, script := range scripts {
				if (executable.Script) == script[pathLength+1:] {
					scriptFound = true
					break
				}
			}
			if !scriptFound {
				return errors.New("Job config invalid: Job '" + job.Id + "' contains a reference to one or more unknown scripts")
			}
		}
	}

	return nil
}

/*
Validate the action configuration
*/
func validateActions(actions []Action, machines []Machine) error {
	for _, action := range actions {
		if action.Id == "" {
			return errors.New("Action config invalid: Each action must have a non-empty id")
		}

		machineFound := false
		for _, machine := range machines {
			if action.Machine == machine.Id || action.Machine == "local" {
				machineFound = true
				break
			}
		}
		if !machineFound {
			return errors.New("Action config invalid: Action '" + action.Id + "' contains a reference to one or more unknown machines")
		}
	}

	return nil
}
