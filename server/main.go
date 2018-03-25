package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/badger"

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
	plexServerInfo, err := db.GetPlexServer()

	if isVerbose && err != nil {
		fmt.Println("plex server:", err.Error())
	}

	if err != nil && err == badger.ErrKeyNotFound || err == badger.ErrEmptyKey {
		// init plex server in datastore
		if err := db.SavePlexServer(datastore.Server{}); err != nil {
			fmt.Printf("failed initializing plex server: %v", err)
		}
	} else if err != nil {
		fmt.Printf("failed retrieving plex server from data store: %v\n", err)
	}

	plexToken, err := db.GetPlexToken()

	if isVerbose && err != nil {
		fmt.Println("plex token error:", err.Error())
	}

	if err != nil && err == badger.ErrKeyNotFound || err == badger.ErrEmptyKey {
		// init plex token in data store
		if err := db.SavePlexServer(datastore.Server{}); err != nil {
			fmt.Printf("failed initializing plex token: %v", err)
		}
	} else if err != nil {
		fmt.Printf("failed retrieving plex token from data store: %v\n", err)
	}

	plexConn, err := plex.New(plexServerInfo.URL, plexToken)

	if err != nil {
		fmt.Printf("failed initializing plex connection: %v\n", err)
	} else if err == nil && isVerbose {
		isPlexInitialized = true
		// no error

		fmt.Println("plex connection initialized")

		isConnected, err := plexConn.Test()

		if err != nil {
			fmt.Printf("testing plex connection failed: %v\n", err)
		}

		fmt.Println("can i connect to my plex server?", isConnected)
	}

	plexConnection := &plexServer{
		plexConn,
		sync.Mutex{},
	}

	router := mux.NewRouter()

	react := make(chan plexUserNotification)

	events := plex.NewNotificationEvents()

	// OnPlaying updates the datastore with user info and status when a plex user starts playing media
	events.OnPlaying(func(n plex.NotificationContainer) {
		currentTime := n.PlaySessionStateNotification[0].ViewOffset
		mediaID := n.PlaySessionStateNotification[0].RatingKey
		sessionID := n.PlaySessionStateNotification[0].SessionKey
		clientID := ""
		userID := "n/a"
		username := "unknown"
		playerType := ""

		var duration int64

		metadata, err := plexConnection.GetMetadata(mediaID)

		if err != nil {
			fmt.Printf("failed to get metadata for key %s: %v\n", mediaID, err)
			return
		}

		duration = metadata.MediaContainer.Metadata[0].Duration
		title := metadata.MediaContainer.Metadata[0].Title

		currentTimeToSeconds := time.Duration(currentTime) * time.Millisecond
		durationToSeconds := time.Duration(duration) * time.Millisecond

		sessions, err := plexConnection.GetSessions()

		if err != nil {
			fmt.Printf("failed to fetch sessions on plex server: %v\n", err)
			return
		}

		for _, session := range sessions.MediaContainer.Video {
			if sessionID != session.SessionKey {
				continue
			}

			userID = session.User.ID
			username = session.User.Title
			sessionID = session.Session.ID
			clientID = session.Player.MachineIdentifier
			playerType = session.Player.Product

			break
		}

		fmt.Printf("%s is playing %s (%s)\n\t%v/%v\n", username, title, mediaID, currentTimeToSeconds, durationToSeconds)

		react <- plexUserNotification{
			ratingKey:   mediaID,
			currentTime: currentTime,
			userID:      userID,
			sessionID:   sessionID,
			clientID:    clientID,
			playerType:  playerType,
		}
	})

	// this anonymous goroutine reads from the datastore and ends plex streams that violate the terms set by one time plex
	go func() {
		if isVerbose {
			fmt.Println("starting goroutine to listen for activity from plex users...")
		}

		for {
			select {
			case r := <-react:
				fmt.Printf("user id: %s, rating key (media id): %s, and current time: %d\n", r.userID, r.ratingKey, r.currentTime)
				userInStore, err := db.GetUser(r.userID)

				if err != nil {
					fmt.Printf("failed to retrieve user from datastore with id of %s\n", r.userID)
					continue
				}

				// end stream
				if userInStore.AssignedMedia.ID != r.ratingKey {
					fmt.Printf("\t%s (%s) is not allowed to watch %s\n\tattempting to terminate stream...\n", userInStore.Name, userInStore.PlexUserID, r.ratingKey)

					if *hasPlexPass {
						// plex pass - will stop playback on whatever player the user is using
						if err := plexConnection.TerminateSession(r.sessionID, "You are not allowed to watch this"); err != nil {
							fmt.Printf("failed to terminate session id %s: %v\n", r.sessionID, err)
						}
					} else {
						fmt.Println("\tclient id:", r.clientID)

						// we cannot use StopPlayback if the user is watching via web browser
						// so we must remove their access to our library
						if r.playerType == "Plex Web" {
							fmt.Println("\tuser is watching via web player -- we must remove their access to library")
							// plexConnection.RemoveFriendAccessToLibrary()
							continue
						}

						// non-plex pass - we can only prevent the user from downloading more data from server
						if plexConnection.StopPlayback(r.clientID); err != nil {
							fmt.Printf("failed to stop playback for user %s (%s), session %s\n", userInStore.Name, r.userID, r.sessionID)
						}
					}

				}
			}
		}
	}()

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
