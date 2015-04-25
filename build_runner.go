package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	protocol  = "http://"
	authority = "localhost:4000"
	rootUrl   = protocol + authority
)

var (
	client = &http.Client{}
)

func buildRunner(queue string, args ...interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in configExtractor", r)
		}
	}()

	defer timeTrack(time.Now(), "buildRunner")
	fmt.Printf("Running build job. Args: %v\n", args)

	projectName := args[0].(string)
	branchName := args[1].(string)
	buildNumber := string(args[2].(json.Number))
	scriptName := args[3].(string)

	url := fmt.Sprintf("%v/projects/%v/branches/%v/builds/%v/config",
		rootUrl, projectName, branchName, buildNumber)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Status code is not 200, it is %v\n", resp.StatusCode)
		return errors.New("Status code != 200")
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	configJson, err := jason.NewObjectFromBytes(body)

	var buffer bytes.Buffer

	index := 0
	scriptMap, _ := configJson.GetObjectArray("scripts")
	for _, script := range scriptMap {
		name, _ := script.GetString("name")
		if name == scriptName {
			break
		}
		index++
	}

	envMap, _ := scriptMap[index].GetObject("ENV")
	for key, value := range envMap.Map() {
		value_str, _ := value.String()
		buffer.WriteString(fmt.Sprintf("export %v=\"%v\"\n", key,
			value_str))
		buffer.WriteString(fmt.Sprintf("echo \"Setting %v to %v\"\n", key,
			value_str))
	}

	commands, _ := scriptMap[index].GetStringArray("commands")
	buffer.WriteString(`
GIT_OUTPUT=$(git clone $RABBIT_CI_REPO --depth=30 2>&1)
if [[ $? -eq 0 ]]; then
    echo $GIT_OUTPUT
else
    echo $GIT_OUTPUT 1>&2
fi
`)
	for _, cmd := range commands {
		buffer.WriteString(fmt.Sprintf("echo '=== Running command: %v'\n",
			cmd))
		buffer.WriteString(fmt.Sprintf("%v\n", cmd))
	}

	dir, err := ioutil.TempDir("", "buildRunner")
	defer os.RemoveAll(dir)

	if err != nil {
		fmt.Printf("Could not create tmp dir. %v\n", err)
		return err
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Printf("Could not change dir. %v\n", err)
		return err
	}

	cl := runCommand("sh", "-c", buffer.String())

	buffer2 := make([]byte, 1024)
	for !cl.finished {
		n, err := cl.pipeReader.Read(buffer2)
		if err != nil {
			cl.pipeReader.Close()
			break
		}

		data := buffer2[0:n]

		sendLog(data, projectName, branchName,
			buildNumber, scriptName)

		for i := 0; i < n; i++ {
			buffer2[i] = 0
		}
	}

	return nil
}

func sendLog(data []byte, projectName, branchName, buildNumber,
	scriptName string) {
	url := fmt.Sprintf("%v%v/projects/%v/branches/%v/builds/%v/log",
		protocol, authority, projectName, branchName, buildNumber)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	values := req.URL.Query()
	values.Add("script", scriptName)
	req.URL.RawQuery = values.Encode()
	req.Header.Set("Content-Type", "text/plain")
	client.Do(req)
}
