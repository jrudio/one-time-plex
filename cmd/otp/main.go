package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/otp-reloaded/http"
	"github.com/jrudio/otp-reloaded/ledis"
)

func main() {
	log.Info(config)
	// open connection to database
	db, err := ledis.Open()

	if err != nil {
		log.WithField("ledis", "open").Fatal(err)
	}

	// create services
	userSrvc := &ledis.UserService{DB: db}

	// attach to http handler
	var h http.Handler
	h.UserService = userSrvc

	endpoints := http.CreateEndpoints(h)

	// start server
	h.ServeHTTP(endpoints, config.Host)
}
