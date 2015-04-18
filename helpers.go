package main

import (
	"time"
	"fmt"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
