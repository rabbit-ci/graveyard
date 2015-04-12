package main

import (
	"fmt"
	"github.com/benmanns/goworker"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func configExtractor(queue string, args ...interface{}) error {
	defer timeTrack(time.Now(), "configExtractor")
	fmt.Printf("Processing job from %s with args: %v\n", queue, args)
	command := exec.Command("worker", "extract-file",
		args[0].(string),
		args[1].(string))

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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func init() {
	goworker.Register("ConfigExtractor", configExtractor)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error: ", err)
	}
}
