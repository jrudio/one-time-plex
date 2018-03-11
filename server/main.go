package main

import (
	"encoding/json"
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

type plexSearchResults struct {
	Title     string `json:"title"`
	Year      string `json:"year"`
	MediaID   string `json:"mediaID"`
	MediaType string `json:"type"`
}

type plexMetadataChildren struct {
	Title     string `json:"title"`
	Year      string `json:"year"`
	MediaID   string `json:"mediaID"`
	MediaType string `json:"type"`
}

type plexFriend struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	// ServerName      string `json:"serverName"`
	// ServerID        string `json:"serverID"`
	// ServerMachineID string `json:"serverMachineID"`
}

type clientResponse struct {
	Result interface{} `json:"result"`
	Err    string      `json:"error"`
}

type restrictedUser struct {
	ID              string `storm:"id" json:"id"`
	Name            string `json:"plexUsername"`
	PlexUserID      string `storm:"unique" json:"plexUserID"`
	AssignedMediaID string `json:"assignedMediaID"`
	Title           string `json:"title"`
}

type usersPayload struct {
	Result []restrictedUser `json:"result"`
}

func (u usersPayload) toBytes() ([]byte, error) {
	return json.Marshal(u)
}

func (r restrictedUser) toBytes() ([]byte, error) {
	return json.Marshal(r)
}

func (c clientResponse) Write(w http.ResponseWriter, errCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "POST GET")

	response, err := json.Marshal(&c)

	if err != nil {
		return err
	}

	if errCode == 0 {
		errCode = http.StatusOK
	}

	w.WriteHeader(errCode)

	_, err = w.Write(response)

	return err
}

// plexServer is a wrapper around plex.Plex to make setting the url and token safe for conccurent use
type plexServer struct {
	*plex.Plex
	lock sync.Mutex
}

func (p *plexServer) setURL(newURL string) {
	p.lock.Lock()
	p.URL = newURL
	p.lock.Unlock()
}

func (p *plexServer) setToken(newToken string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Token = newToken
}

func (p *plexServer) getURL() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.URL
}

type plexUserNotification struct {
	// ratingKey aka media id
	ratingKey   string
	currentTime int64
	userID      string
	sessionID   string
	clientID    string
	playerType  string
}

// AddUser adds a user that needs to be monitored
func AddUser(db datastore.Store, plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user datastore.User
		var resp clientResponse

		// check for required parameters
		if err := r.ParseForm(); err != nil {
			resp.Err = "failed to parse form"
			resp.Write(w, http.StatusBadRequest)
			return
		}

		plexUserID := r.FormValue("plexUserID")
		mediaID := r.FormValue("mediaID")

		user.PlexUserID = plexUserID
		user.AssignedMedia.ID = mediaID
		user.AssignedMedia.Status = "not started"

		if user.PlexUserID == "" || user.AssignedMedia.ID == "" {
			resp.Err = "missing 'plexuserid' and/or 'mediaID' in the post form body"
			resp.Write(w, http.StatusBadRequest)
			return
		}

		metadata, err := plexConnection.GetMetadata(user.AssignedMedia.ID)

		if err != nil {
			resp.Err = fmt.Sprintf("failed to fetch title for media id %s: %v", user.AssignedMedia.ID, err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		metadataLen := len(metadata.MediaContainer.Metadata)
		// metadataLen, err := strconv.Atoi(metadata.Size)

		// if err != nil {
		// 	resp.Err = fmt.Sprintf("(%s) failed to convert metadata length to int: %v", user.AssignedMedia.ID, err)
		// 	resp.Write(w, http.StatusInternalServerError)
		// 	return
		// }

		if metadataLen > 0 {
			data := metadata.MediaContainer.Metadata[0]

			var title string

			if isVerbose {
				fmt.Printf("metadata for %s: title: %s, grandparent title: %s, parent title: %s, type: %s\n",
					data.RatingKey,
					data.Title,
					data.GrandparentTitle,
					data.ParentTitle,
					data.Type,
				)
			}

			// only handle episodes or movies -- season and show return error
			if data.Type != "episode" && data.Type != "movie" {
				resp.Err = "assigned media must be of type 'movie' or 'episode' - received type: '" + data.Type + "'"
				resp.Write(w, http.StatusBadRequest)
				return
			}

			// combine show name, season, and episode name if type == episode
			if data.Type == "episode" {
				title = data.GrandparentTitle + ": " + data.ParentTitle + " - " + data.Title
			} else {
				// type is movie just need the title
				title = metadata.MediaContainer.Metadata[0].Title
			}

			user.Title = title
		}

		// we can assume the user being added is a plex friend
		friends, err := plexConnection.GetFriends()

		if err != nil {
			resp.Err = fmt.Sprintf("failed to fetch plex friends: %v", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		friendUsername, err := getFriendUsernameByID(friends, user.PlexUserID)

		if err != nil {
			resp.Err = fmt.Sprintf("failed to get plex username: %v", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		user.Name = friendUsername

		if err := db.SaveUser(user); err != nil {
			resp.Err = fmt.Sprintf("failed to save user: %v\n", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		resp.Result = user

		resp.Write(w, http.StatusOK)
	}
}

func getFriendUsernameByID(friends []plex.Friends, id string) (string, error) {
	for _, friend := range friends {
		friendID := strconv.Itoa(friend.ID)

		if friendID == id {
			return friend.Username, nil
		}
	}

	return "", fmt.Errorf("could not find friend with that id")
}

// GetAllUsers returns all monitored users
func GetAllUsers(db datastore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var resp clientResponse

		users, err := db.GetAllUsers()

		if err != nil {
			resp.Err = fmt.Sprintf("failed to retrieve users: %v\n", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		fmt.Println(users)

		resp.Result = users
		resp.Write(w, http.StatusOK)
	}
}

// RemoveAllUsers clears out the monitored users
func RemoveAllUsers(db datastore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := clientResponse{}

		users, err := db.GetAllUsers()

		if err != nil {
			resp.Err = fmt.Sprintf("failed to fetch all user ids: %v", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		userIDs := make([]string, len(users))

		i := 0

		for userID := range users {
			userIDs[i] = userID

			i++
		}

		if err := db.DeleteUsers(userIDs); err != nil {
			resp.Err = fmt.Sprintf("failed to delete users: %v", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		resp.Result = true
		resp.Write(w, http.StatusOK)
	}
}

// RemoveUser removes a monitored user
func RemoveUser(db datastore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVars := mux.Vars(r)
		resp := clientResponse{}

		userID, ok := pathVars["id"]

		if !ok {
			resp.Result = "one time plex failed to parse id -- developer error"
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		if userID == "" {
			resp.Err = "user id is required"
			resp.Write(w, http.StatusBadRequest)
			return
		}

		if err := db.DeleteUser(userID); err != nil {
			resp.Err = fmt.Sprintf("failed to delete user %s: %v", userID, err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		resp.Result = true
		resp.Write(w, http.StatusOK)
	}
}

// GetMetadataFromPlex fetches the metadata of plex media
func GetMetadataFromPlex(plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := clientResponse{}

		query := r.URL.Query()

		mediaID := query.Get("mediaid")

		if mediaID == "" {
			response.Err = "mediaid required"
			response.Write(w, http.StatusBadRequest)
			return
		}

		metadata, err := plexConnection.GetMetadataChildren(mediaID)

		if err != nil {
			response.Err = fmt.Sprintf("failed plexMetadataChildrento fetch metadata for %s: %v", mediaID, err)
			response.Write(w, http.StatusInternalServerError)
			return
		}

		response.Result = filterMetadata(metadata)
		response.Write(w, http.StatusOK)
	}
}

func filterMetadata(results plex.MetadataChildren) []plexMetadataChildren {
	resultsLen := len(results.MediaContainer.Metadata)
	filteredMetadata := make([]plexMetadataChildren, resultsLen)

	for i := 0; i < resultsLen; i++ {
		metadata := results.MediaContainer.Metadata[i]

		filteredMetadata[i] = plexMetadataChildren{
			Title:     metadata.Title,
			Year:      "N/A",
			MediaID:   metadata.Key,
			MediaType: metadata.Type,
		}
	}

	return filteredMetadata
}

// GetPlexFriends fetches plex friends associated with the set plex token
func GetPlexFriends(plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := clientResponse{}

		friends, err := plexConnection.GetFriends()

		if err != nil {
			response.Err = fmt.Sprintf("failed to fetch plex friends: %v", err)
			response.Write(w, http.StatusInternalServerError)
			return
		}

		response.Result = filterFriends(friends)
		response.Write(w, http.StatusOK)
	}
}

func filterFriends(friends []plex.Friends) []plexFriend {
	friendsLen := len(friends)
	filteredFriends := make([]plexFriend, friendsLen)

	for i := 0; i < friendsLen; i++ {
		friend := friends[i]

		filteredFriends[i] = plexFriend{
			ID:       strconv.Itoa(friend.ID),
			Username: friend.Title,
		}
	}

	return filteredFriends
}

func filterSearchResults(results plex.SearchResults) []plexSearchResults {
	var newResults []plexSearchResults

	count := results.MediaContainer.Size

	if count == 0 {
		return newResults
	}

	for _, r := range results.MediaContainer.Metadata {
		filtered := plexSearchResults{
			MediaType: r.Type,
			MediaID:   r.RatingKey,
			Title:     r.Title,
			Year:      "N/A", // default to n/a if we can't convert to string
		}

		if year := strconv.FormatInt(r.Year, 10); year != "" {
			filtered.Year = year
		}

		newResults = append(newResults, filtered)
	}

	return newResults
}

// SearchPlex is an endpoint that will search your Plex Media Server for media
func SearchPlex(plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var resp clientResponse

		searchQuery := r.URL.Query().Get("title")

		if searchQuery == "" {
			resp.Err = "missing search query: 'title'"
			resp.Write(w, http.StatusBadRequest)
			return
		}

		results, err := plexConnection.Search(searchQuery)

		if err != nil {
			resp.Err = fmt.Sprintf("search on plex media server failed: %v", err)
			resp.Write(w, http.StatusInternalServerError)
			return
		}

		// filter results with relevant information
		resp.Result = filterSearchResults(results)

		resp.Write(w, http.StatusOK)
	}
}

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

// ConfigurePlexServer retrieve, add, edit, and delete plex server from one time plex
func ConfigurePlexServer(db datastore.Store, plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var resp clientResponse

		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				resp.Err = fmt.Sprintf("parsing body failed: %v", err)
				resp.Write(w, http.StatusBadRequest)
				return
			}

			plexURL := r.PostFormValue("url")
			plexToken := r.PostFormValue("token")

			resp.Result = fmt.Sprintf("received %s, %s", plexURL, plexToken)
		case "PUT":
			if err := r.ParseForm(); err != nil {
				resp.Err = fmt.Sprintf("parsing body failed: %v", err)
				resp.Write(w, http.StatusBadRequest)
				return
			}

			if plexURL := r.FormValue("url"); plexURL != "" {
				plexConnection.setURL(plexURL)

				if err := db.SavePlexServer(datastore.Server{
					Name: "default",
					URL:  plexURL,
				}); err != nil {
					fmt.Printf("saving plex url failed: %v", err)
				}
			}
			if plexToken := r.FormValue("token"); plexToken != "" {
				plexConnection.setToken(plexToken)

				if err := db.SavePlexToken(plexToken); err != nil {
					fmt.Printf("saving plex token failed: %v", err)
				}
			}

			resp.Result = true

		default:
			resp.Err = "unknown method"
			resp.Write(w, http.StatusMethodNotAllowed)
			return
		}

		resp.Write(w, http.StatusOK)
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

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/plex/server", ConfigurePlexServer(db, plexConnection)).Methods("POST", "PUT", "DELETE")

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
