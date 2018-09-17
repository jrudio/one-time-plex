package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jrudio/go-plex-client"
	"github.com/jrudio/one-time-plex/server/datastore"
	homedir "github.com/mitchellh/go-homedir"
)

const databaseFileName = "onetimeplex.db"

var (
	isVerbose         bool
	isPlexInitialized bool
)

func cleanup(ctrlC chan os.Signal, shutdown chan bool, db datastore.Store) {
	for {
		select {
		case <-ctrlC:
			// funnel the interrupt to "shutdown"
			shutdown <- true
		case s := <-shutdown:
			if !s {
				continue
			}

			if isVerbose {
				fmt.Println("gracefully closing data store...")
			}

			fmt.Println("shutting down...")

			db.Close()
			os.Exit(0)
		}
	}
}

// func startPlexServices() {}

func main() {
	// grab optional params
	listenOnAddress := flag.String("host", ":6969", "where the server should listen on. Default is localhost:6969")
	isVerboseFlag := flag.Bool("verbose", false, "spit out more information")
	hasPlexPass := flag.Bool("plexpass", false, "utilize plex pass features")

	removeLockFile := flag.Bool("unlock", false, "remove the lock file - Do this if one time plex crashes and you cannot start again")

	flag.Parse()

	if *isVerboseFlag {
		isVerbose = *isVerboseFlag
	}

	// create channels to shutdown gracefully
	ctrlC := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)

	signal.Notify(ctrlC, os.Interrupt)

	// get user home directory to house our datastore
	storeDirectory, err := homedir.Dir()

	if err != nil {
		fmt.Printf("could not find home directory: %v\n", err)
		os.Exit(1)
	}

	storeDirectory = filepath.Join(storeDirectory, ".one-time-plex")

	// 'unlock' flag was used
	if *removeLockFile {
		lockFile := filepath.Join(storeDirectory, "LOCK")

		if err := os.Remove(lockFile); err != nil {
			fmt.Printf("failed to remove lock file: %v\n", err)
		} else {
			fmt.Println("removed lock file")
		}

		os.Exit(1)
	}

	// connect to datastore
	db, err := datastore.InitDataStore(storeDirectory, isVerbose)

	// handle interrupts - i.e. control-c
	go cleanup(ctrlC, shutdown, db)

	if err != nil {
		fmt.Printf("failed to open or create the datastore: %v\n", err)
		shutdown <- true
		return
	}

	// check if app secret exists in datastore
	appSecret := []byte("iAmAseCReTuSEdTOENcrYp")

	if secret := db.GetSecret(); len(secret) != 0 {
		db.Secret = secret
	} else {
		// if not, create new appsecret and save
		// append a random number to make appsecret unique
		appSecret = append(appSecret, []byte(strconv.FormatInt(time.Now().Unix(), 10))...)

		if err := db.SaveSecret(appSecret); err != nil {
			fmt.Printf("failed to save app secret: %v\n", err)
			shutdown <- true
			return
		}

		db.Secret = appSecret
	}

	// init plex stuff
	plexConnection := initPlex(db)

	router := mux.NewRouter()

	react := make(chan plexUserNotification)

	events := plex.NewNotificationEvents()

	// OnPlaying updates the datastore with user info and status when a plex user starts playing media
	events.OnPlaying(OnPlaying(react, plexConnection))

	// this goroutine reads from the datastore and ends plex streams that violate the terms set by one time plex
	go streamGuard(react, db, plexConnection, *hasPlexPass)

	if isVerbose {
		fmt.Println("subscribing to plex server notifications...")
	}

	plexConnection.SubscribeToNotifications(events, ctrlC, func(err error) {
		fmt.Printf("failed to subscribe to server: %v\n", err)
		// shutdown <- true
		// return
	})

	fmt.Println("connected to plex server via websocket\nlistening for events...")

	if isVerbose {
		fmt.Println("after subscribing")
	}

	// serve front end
	router.HandleFunc("/", ServeAsset("build/index.html"))
	router.PathPrefix("/static/").Handler(http.FileServer(assetFS()))
	router.HandleFunc("/favicon.ico", ServeAsset("build/favicon.ico"))

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/plex/server", ConfigurePlexServer(db, plexConnection)).Methods("POST", "PUT", "DELETE", "GET", "OPTIONS")
	apiRouter.HandleFunc("/plex/test", TestPlexConnection).Methods("GET")
	apiRouter.HandleFunc("/plex/pin", GetPlexPin(db)).Methods("GET")

	// add new one-time user
	apiRouter.HandleFunc("/users/add", AddUser(db, plexConnection)).Methods("POST")

	// list one-time users
	apiRouter.HandleFunc("/users", GetAllUsers(db)).Methods("GET")

	// clear all one-time users
	apiRouter.HandleFunc("/users/clear", RemoveAllUsers(db)).Methods("DELETE")

	apiRouter.HandleFunc("/user/{id}", RemoveUser(db)).Methods("DELETE", "OPTIONS")

	// search media on plex
	apiRouter.HandleFunc("/search", SearchPlex(plexConnection)).Methods("GET")

	// get plex friends
	apiRouter.HandleFunc("/friends", GetPlexFriends(plexConnection)).Methods("GET")

	// get child data from plex
	apiRouter.HandleFunc("/metadata", GetMetadataFromPlex(plexConnection)).Methods("GET")

	// fmt.Printf("serving one time plex on %s\n", config.Host)

	fmt.Println("server listening on " + *listenOnAddress)

	if err := http.ListenAndServe(*listenOnAddress, router); err != nil {
		fmt.Printf("server failed to start: %v\n", err)
		shutdown <- true
		return
	}
}
