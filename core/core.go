/*
The public interface to the core part of the application shared by the command
line interface as well as the CI server
*/

package core

/*
Get a new Api for executing commands locally
*/
func NewLocalApi(path string) (Api, error) {
	return LocalApi{path: path}, nil
}

/*
Get a new Api for executing commands on the remote server
*/
func NewRemoteApi(path string) (Api, error) {
	return RemoteApi{path: path}, nil
}
