package common

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func ErrorHandler(err error, fatal bool) {
	ts := getTS()
	if err != nil {
		log.Errorf("error: %v %v", err, ts)

		if fatal {
			panic(err)
		}
	}
}

func LogInfo(info string) {
	ts := getTS()
	log.Infof("%v, %v", info, ts)
}

func getTS() string {
	currentTime := time.Now()
	// Format the current time as a timestamp
	timestamp := currentTime.Format(time.RFC3339)
	return timestamp
}
