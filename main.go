package main

import (
	log "github.com/sirupsen/logrus"
	"goals/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "app.Run",
			"error":    err,
		}).Error("Unable to start app")
		panic(err)
	}
}
