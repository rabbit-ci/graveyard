package main

import (
	"fmt"
	"github.com/benmanns/goworker"
	"os/exec"
)

func myFunc(queue string, args ...interface{}) error {
	out, err := exec.Command("worker", "extract-file",
		"git@github.com:hunterboerner/rabbitci-test-repo.git",
		"ebd9c8282f992e009e04dcbf6c1e7df07f554c48", "--file=.rabbitci.json").Output()
	fmt.Println(err)
	fmt.Println(string(out))
	fmt.Printf("From %s, %v\n", queue, args)
	return nil
}

func init() {
	goworker.Register("MyClass", myFunc)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error: ", err)
	}
}
