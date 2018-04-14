package main

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger"
	plex "github.com/jrudio/go-plex-client"
	"github.com/jrudio/one-time-plex/server/datastore"
)

func initPlex(db datastore.Store) *plexServer {
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

	return &plexServer{
		plexConn,
		sync.Mutex{},
	}
}
