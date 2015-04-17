package worker

import (
	"fmt"
	"time"
	"os/exec"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
)

func configExtractor(queue string, args ...interface{}) error {
	defer timeTrack(time.Now(), "configExtractor")

	repo := args[0].(string)
	commit := args[1].(string)
	branch := args[2].(string)

	fmt.Printf("Processing job from %s with args: %v\n", queue, args)
	command := exec.Command("worker", "extract-file",
		repo, commit)

	command.Env = os.Environ()
	command.Env = append(command.Env, "RUST_BACKTRACE=1")

	out, err := command.CombinedOutput()

	if err != nil {
		// TODO: We need to tell the server that we encountered an error and give it the backtrace.
		fmt.Println(string(out))
		fmt.Println("File extraction errored!")
		return err
	}

	url := "http://localhost:4000/config_extraction" // TODO: This needs to be configurable
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(out))

	values := req.URL.Query()
	values.Add("commit", commit)
	values.Add("branch", branch)
	values.Add("repo", repo)
	req.URL.RawQuery = values.Encode()

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	fmt.Printf("Finished processing job. Args: %v\n", args)
	return nil
}

