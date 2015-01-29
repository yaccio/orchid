/*
Definition of and methods for logs including log persistence
*/

package core

import (
	"encoding/json"
	"github.com/dchest/uniuri"
	"io/ioutil"
	"os"
	"time"
)

/*
Definition of the log type
*/
type Log struct {
	Id        string
	JobId     string
	Status    string
	StartTime time.Time
	EndTime   time.Time
}

/*
Save the log to the logs configuration file
*/
func (l Log) save(path string) error {
	// Keep the file open as a write lock, ensuring that no lost updates occur
	f, err := os.OpenFile(path+"/logs.json", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err == nil {
		// New file. Write empty list to it
		_, err = f.WriteString("[]")
		if err != nil {
			f.Close()
			return err
		}
	}
	defer f.Close()

	logs, err := loadLogs(path)
	if err != nil {
		return err
	}

	found := false
	for i, log := range logs {
		if l.Id == log.Id {
			found = true
			logs[i] = l
			break
		}
	}

	if !found {
		logs = append(logs, l)
	}

	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path+"/logs.json", data, 0644)
}

/*
Indicate that the log has started, setting the start time and updating the
persistent log configuration
*/
func (l Log) start(path string) (Log, error) {
	l.StartTime = time.Now()
	l.Status = "Started"
	return l, l.save(path)
}

/*
Indicate that the log has finished, setting the end time and updating the
persistent log configuration
*/
func (l Log) finish(path string, file *os.File) (Log, error) {
	l.EndTime = time.Now()
	l.Status = "Finished"
	return l, l.saveAndWriteToLog(path, file, "Finished")
}

/*
Indicate that the log has encountered and error, setting the end time and
updating the persistent log configuration
*/
func (l Log) error(path string, file *os.File) (Log, error) {
	l.EndTime = time.Now()
	l.Status = "Error"
	return l, l.saveAndWriteToLog(path, file, "Error")
}

/*
Helper method for saving the log and writing a terminating line to the log
output file
*/
func (l Log) saveAndWriteToLog(path string, file *os.File, text string) error {
	err := l.save(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString("-----" + text + "-----")
	return err
}

/*
Create a new log, assigning it a new identifier
*/
func newLog(jobId string) Log {
	return Log{
		Id:     uniuri.New(),
		JobId:  jobId,
		Status: "New",
	}
}

/*
Load all logs stored locally
*/
func loadLogs(path string) ([]Log, error) {
	// Load the file
	logs := &[]Log{}
	data, err := ioutil.ReadFile(path + "/logs.json")
	if err != nil {
		return []Log{}, err
	}

	err = json.Unmarshal(data, &logs)
	if err != nil {
		return []Log{}, err
	}

	return *logs, nil
}
