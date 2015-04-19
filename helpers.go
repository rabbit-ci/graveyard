package main

import (
	"fmt"
	"github.com/rabbit-ci/logstreamer"
	"io"
	"log"
	"os/exec"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

type CommandLog struct {
	pipeReader *io.PipeReader
	finished   bool
}

func runCommand(name string, arg ...string) *CommandLog {
	pipeReader, pipeWriter := io.Pipe()
	logger := log.New(pipeWriter, "--> ", log.Ldate|log.Ltime)

	logStreamerOut := logstreamer.NewLogstreamer(logger, "stdout", false)
	logStreamerErr := logstreamer.NewLogstreamer(logger, "stderr", true)

	cmd := exec.Command(name, arg...)
	cmd.Stderr = logStreamerErr
	cmd.Stdout = logStreamerOut

	// Reset any error we recorded
	logStreamerErr.FlushRecord()

	cl := CommandLog{pipeReader: pipeReader, finished: false}
	go func() {
		err := cmd.Start()
		// Failed to spawn?
		if err != nil {
			logger.Printf("ERROR could not spawn command. %s \n", err.Error())
		}

		errf := cmd.Wait()
		if errf != nil {
			logger.Printf("Someting went wrong. %s \n", err.Error())
		}
		cl.finished = true
		logger.Print("Finished")
	}()

	return &cl
}
