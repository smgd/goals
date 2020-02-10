package main

import (
	"goals/app/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := server.Run(); err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "app.Run",
			"error":    err,
		}).Error("Unable to start app")
		panic(err)
	}
}
