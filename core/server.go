/*
Definition of and methods for loading and validating the remote CI server
configuration files
*/

package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

/*
Type defining the remote server setup configuration
*/
type Server struct {
	Address string
	Port    string
	Secret  string
}

/*
Load the remote server configuration
*/
func loadServer(path string) (Server, error) {
	server := &Server{}
	data, err := ioutil.ReadFile(path + "/server.json")
	if err != nil {
		return Server{}, err
	}

	err = json.Unmarshal(data, &server)
	if err != nil {
		return Server{}, err
	}

	err = validateServer(*server)
	if err != nil {
		return Server{}, err
	}

	return *server, nil
}

/*
Validate the remote server configuration
*/
func validateServer(server Server) error {
	if server.Address == "" {
		return errors.New("Server config invalid: Server must have a non-empty Address")
	}
	if server.Port == "" {
		return errors.New("Server config invalid: Server must have a non-empty Port")
	}

	return nil
}
