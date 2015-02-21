package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
Application entry point
Reads a line from stdin and passes it as arguments to the orchid application
*/
func main() {
	buffer := bufio.NewReader(os.Stdin)
	line, _, err := buffer.ReadLine()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	os.Setenv("PATH", "/bin:/usr/bin")

	args := strings.Fields(string(line))

	cmd := exec.Command("orchid", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
}
