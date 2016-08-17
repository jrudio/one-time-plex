package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/go-plex-client"
	"github.com/jrudio/otp-reloaded/http"
	"github.com/jrudio/otp-reloaded/ledis"
	"github.com/jrudio/otp-reloaded/monitor"
)

func main() {
	// open connection to database
	db, err := ledis.Open()

	if err != nil {
		log.WithField("ledis", "open").Fatal(err)
	}

	// create services
	var plexConn *plex.Plex
	plexConn, err = plex.New(config.Plex.Host, config.Plex.Token)

	if err != nil {
		log.WithField("plex", "new").Fatal(err)
	}

	userSrvc := &ledis.UserService{DB: db}
	plexMonitorSrvc := &monitor.PlexMonitorService{
		DB:   db,
		Plex: plexConn,
	}

	// attach to http handler
	var h http.Handler
	h.UserService = userSrvc
	h.PlexMonitorService = plexMonitorSrvc
	h.Plex = plexConn

	endpoints := http.CreateEndpoints(h)

	// start server
	log.Info("serving One Time Plex at " + config.Host)
	h.ServeHTTP(endpoints, config.Host)
}
