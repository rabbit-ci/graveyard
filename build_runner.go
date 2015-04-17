package worker

import (
	"time"
)

func buildRunner(queue string, args ...interface{}) error {
	defer timeTrack(time.Now(), "buildRunner")
	return nil
}
