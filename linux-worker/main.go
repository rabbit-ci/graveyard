package main

import (
	"fmt"
	"github.com/benmanns/goworker"
)

func init() {
	goworker.Register("ConfigExtractor", configExtractor)
	goworker.Register("BuildRunner", buildRunner)
}

func main() {
	if err := goworker.Work(); err != nil {
		// This doesn't appear to be doing anything
		fmt.Println("Error: ", err)
	}
}
