/*
Implementation of the Api for executing commands remotely
*/

package core

import (
	"encoding/json"
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

	var log Log
	err = json.Unmarshal(result, &log)
	if err != nil {
		return Log{}, err
	}

	return log, nil
}

/*
List all jobs stored remotely
*/
func (a RemoteApi) ListLogs() ([]Log, error) {
	url := fmt.Sprintf("/logs/list")
	result, err := a.doRequest("GET", url)
	if err != nil {
		return []Log{}, err
	}

	var logs []Log
	err = json.Unmarshal(result, &logs)
	if err != nil {
		return []Log{}, err
	}

	return logs, nil
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
func (a RemoteApi) doRequest(method, url string) ([]byte, error) {
	server, err := loadServer(a.path)
	if err != nil {
		return []byte{}, err
	}

	url = fmt.Sprintf("%s:%s%s", server.Address, server.Port, url)

	req, requestErr := http.NewRequest(method, url, nil)
	if requestErr != nil {
		return []byte{}, requestErr
	}
	if server.Secret != "" {
		req.Header.Set("X-Orchid-Secret", server.Secret)
	}

	client := &http.Client{}
	response, responseErr := client.Do(req)
	if responseErr != nil {
		return []byte{}, responseErr
	}
	defer response.Body.Close()

	body, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return []byte{}, bodyErr
	}

	return body, nil
}
