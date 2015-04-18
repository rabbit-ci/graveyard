package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
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

	protocol := "http://"
	authority := "localhost:4000"
	url := fmt.Sprintf("%v%v/projects/%v/branches/%v/builds/%v/config_file",
		protocol, authority, projectName, branchName, buildNumber)
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
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

	var configJson map[string]interface{}
	if err := json.Unmarshal(body, &configJson); err != nil {
		return err
	}

	var buffer bytes.Buffer
	for key, value := range configJson["ENV"].(map[string]interface{}) {
		buffer.WriteString(fmt.Sprintf("export %v=\"%v\"\n", key,
			value))
	}

	for _, cmd := range configJson["commands"].([]interface{}) {
		buffer.WriteString(fmt.Sprintf("echo '=== Running command: %v'\n",
			cmd.(string)))
		buffer.WriteString(fmt.Sprintf("%v\n", cmd.(string)))
	}

	dir, err := ioutil.TempDir("", "buildRunner")
	fmt.Println(dir)
	defer os.RemoveAll(dir)

	if err != nil {
		fmt.Printf("Could not create tmp dir. %v\n", err)
		return err
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Printf("Could not change dir. %v\n", err)
		return err
	}

	out, err := exec.Command("sh", "-c", buffer.String()).Output()
	fmt.Println(string(out))

	return nil
}
