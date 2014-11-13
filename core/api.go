/*
Definition of the interface defining the available commands to execute both
locally and remotely
*/

package core

import (
	"io"
)

/*
The interface definition
*/
type Api interface {
	//Run the job with the given id
	RunJob(jobId string) (Log, error)

	//List all stored logs
	ListLogs() ([]Log, error)

	// Get the output stored in the log with the given id
	GetLogOutput(logId string, writer io.Writer) error
}
