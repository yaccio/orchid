/*
Implementation of the Api for executing commands locally
*/

package core

import (
	"fmt"
	"github.com/ActiveState/tail"
	"io"
)

/*
The local API type
*/
type LocalApi struct {
	path string
}

/*
Execute the job with the given id locally
*/
func (a LocalApi) RunJob(jobId string) (Log, error) {
	log := newLog(jobId)

	pipeline, err := buildPipeline(a.path, jobId, log)
	if err != nil {
		return Log{}, err
	}

	go func() {
		pipeline.Run(a.path)
	}()

	return log, nil
}

/*
List all existing logs stored locally
*/
func (a LocalApi) ListLogs() ([]Log, error) {
	return loadLogs(a.path)
}

/*
Get the output stored locally in the log with the given id
*/
func (a LocalApi) GetLogOutput(logId string, writer io.Writer) error {
	t, err := tail.TailFile(a.path+"/logs/"+logId, tail.Config{Follow: true})
	if err != nil {
		return err
	}

	for line := range t.Lines {
		if line.Text == "-----Finished-----" || line.Text == "-----Error-----" {
			break
		}
		fmt.Fprintln(writer, line.Text)
	}

	return nil
}
