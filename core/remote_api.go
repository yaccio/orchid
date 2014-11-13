/*
Implementation of the Api for executing commands remotely
*/

package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

/*
The remote API type
*/
type RemoteApi struct {
	path string
}

/*
Execute the job with the given id remotely
*/
func (a RemoteApi) RunJob(jobId string) (Log, error) {
	url := fmt.Sprintf("/runs/%s", jobId)
	result, err := a.doRequest("POST", url)
	if err != nil {
		return Log{}, err
	}

	//TODO parse as Log
	return Log{Id: result}, nil
}

/*
List all jobs stored remotely
*/
func (a RemoteApi) ListLogs() ([]Log, error) {
	//TODO
	return []Log{}, nil
}

/*
Get the output stored remotely in the log with the given id
*/
func (a RemoteApi) GetLogOutput(logId string, writer io.Writer) error {
	//TODO
	return nil
}

/*
Helper method for executing a HTTP request on the remote server as defined in
the configurations
*/
func (a RemoteApi) doRequest(method, url string) (string, error) {
	server, err := loadServer(a.path)
	if err != nil {
		return "", err
	}

	url = fmt.Sprintf("%s:%s%s", server.Address, server.Port, url)

	req, requestErr := http.NewRequest(method, url, nil)
	if requestErr != nil {
		return "", requestErr
	}
	if server.Secret != "" {
		req.Header.Set("X-Orchid-Secret", server.Secret)
	}

	client := &http.Client{}
	response, responseErr := client.Do(req)
	if responseErr != nil {
		return "", responseErr
	}
	defer response.Body.Close()

	body, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return "", bodyErr
	}

	return string(body), nil
}
