package main

import (
	"fmt"
	"testing"
)

func TestRunCommand(t *testing.T) {
	cl := runCommand("sh", "-c", "echo 'wut'; sleep 2; echo 'hey'; echo 'wat'; sleep 2; echo 'finish'")
	buffer := make([]byte, 1024)
	for !cl.finished {
		n, err := cl.pipeReader.Read(buffer)
		if err != nil {
			cl.pipeReader.Close()
			break
		}

		data := buffer[0:n]
		fmt.Print(string(data))
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}
