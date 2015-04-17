package worker

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
		fmt.Println("Error: ", err)
	}
}
