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
	buffer   *io.PipeReader
	finished bool
}

func runCommand(name string, arg ...string) *CommandLog {
	// var b bytes.Buffer
	pipeReader, pipeWriter := io.Pipe()
	// writer := bufio.NewWriter(&b)
	logger := log.New(pipeWriter, "--> ", log.Ldate|log.Ltime)

	// Setup a streamer that we'll pipe cmd.Stdout to
	logStreamerOut := logstreamer.NewLogstreamer(logger, "stdout", false)
	// Setup a streamer that we'll pipe cmd.Stderr to.
	// We want to record/buffer anything that's written to this (3rd argument true)
	logStreamerErr := logstreamer.NewLogstreamer(logger, "stderr", true)

	cmd := exec.Command(name, arg...)
	cmd.Stderr = logStreamerErr
	cmd.Stdout = logStreamerOut

	// Reset any error we recorded
	logStreamerErr.FlushRecord()

	cl := CommandLog{buffer: pipeReader, finished: false}
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
