package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func logmessage(t time.Time) {
	log.Info("Hello from the Logsearch-for-CloudFoundry")
}

func main() {
	doEvery(10*time.Millisecond, logmessage)
}
