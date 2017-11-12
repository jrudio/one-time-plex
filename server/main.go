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

	"github.com/asdine/storm"
	"github.com/gorilla/mux"
	"github.com/jrudio/go-plex-client"
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

type plexFriend struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	ServerName      string `json:"serverName"`
	ServerID        string `json:"serverID"`
	ServerMachineID string `json:"serverMachineID"`
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

type server struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (s server) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func unserializeServer(serializedServer []byte) (server, error) {
	var s server

	err := json.Unmarshal(serializedServer, &s)

	return s, err
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

// OnPlay reads and reacts to the webhook sent from Plex
func OnPlay(db *storm.DB, plexConnection *plex.Plex) func(wh plex.Webhook) {
	return func(wh plex.Webhook) {
		userID := strconv.Itoa(wh.Account.ID)
		username := wh.Account.Title
		title := wh.Metadata.Title
		mediaID := wh.Metadata.RatingKey

		fmt.Printf("%s (%s) has started playing %s (%s)\n", username, userID, title, mediaID)

		// is this a user we need to check?
		var user restrictedUser

		if err := db.One("PlexUserID", userID, &user); err != nil {
			// user not in database so we don't care about them
			fmt.Printf("user %s (%s) is not in database\n", username, userID)
			return
		}

		if user.AssignedMediaID == mediaID {
			// user is watching what they were assigned
			fmt.Printf("user %s (%s) is watching %s which is ok\n", username, userID, mediaID)
			return
		}

		// Obtain session id
		//
		// We will assume the plexConnection is the server that sent this webhook
		sessions, err := plexConnection.GetSessions()

		if err != nil {
			fmt.Printf("not terminating user: %s (%s) \n\tfailed to grab sessions from plex server: %v\n", username, userID, err)
			return
		}

		var sessionID string

		for _, session := range sessions.MediaContainer.Video {
			if session.User.ID != userID {
				continue
			}

			sessionID = session.Session.ID
			break
		}

		// kill session
		fmt.Printf("Terminating %s (%s)'s session as they are not supposed to be watching %s (%s)\n", username, userID, title, mediaID)
		plexConnection.TerminateSession(sessionID, "One Time Plex: You are not allowed to watch that")
		// fmt.Printf("%d is now playing: %s (%s)\n", userID, title, mediaID)
	}
}

// OnStop will stop monitoring the user and unshare the Plex library
func OnStop(db *storm.DB, plexConnection *plex.Plex) func(wh plex.Webhook) {
	return func(wh plex.Webhook) {
		// remove from our database
		username := wh.Account.Title
		userID := wh.Account.ID

		var user restrictedUser

		if err := db.One("PlexUserID", userID, &user); err != nil {
			// user not in database, don't care
			fmt.Printf("user %s (%d) is not in database\n", username, userID)
			return
		}

		if err := db.DeleteStruct(&user); err != nil {
			fmt.Printf("user %s (%d) removal failed\n", username, userID)
			return
		}

		// unshare the Plex library
		_, err := plexConnection.RemoveFriend(strconv.Itoa(userID))

		if err != nil {
			fmt.Printf("failed to unshare library with %s (%d)\n\tplease remove them manually\n", username, userID)
			return
		}

		fmt.Printf("%s (%d) has finished viewing: %s", username, userID, wh.Metadata.Title)
	}
}

// AddUser adds a user that needs to be monitored
// func AddUser(db *storm.DB, plexConnection *plex.Plex) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user restrictedUser
// 		var resp clientResponse

// 		// check for required parameters
// 		defer r.Body.Close()

// 		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 			resp.Err = fmt.Sprintf("failed to decode json body: %v", err)
// 			resp.Write(w)
// 			return
// 		}

// 		if user.PlexUserID == "" || user.AssignedMediaID == "" {
// 			resp.Err = "missing 'plexuserid' and/or 'mediaID' in the post form body"
// 			resp.Write(w)
// 			return
// 		}

// 		user.ID = xid.New().String()
// 		// user.Name = plexUsername
// 		// user.PlexUserID = plexUserID
// 		// user.AssignedMediaID = mediaID

// 		metadata, err := plexConnection.GetMetadata(user.AssignedMediaID)

// 		if err != nil {
// 			resp.Err = fmt.Sprintf("failed to fetch title for media id %s: %v", user.AssignedMediaID, err)
// 			resp.Write(w)
// 			return
// 		}

// 		metadataLen, err := strconv.Atoi(metadata.Size)

// 		if err != nil {
// 			resp.Err = fmt.Sprintf("(%s) failed to convert metadata length to int: %v", user.AssignedMediaID, err)
// 			resp.Write(w)
// 			return
// 		}

// 		if metadataLen > 0 {
// 			data := metadata.Metadata[0]

// 			var title string

// 			// combine show name, season, and episode name if type == episode
// 			if data.Type == "episode" {
// 				title = data.GrandparentTitle + ": " + data.ParentTitle + " - " + data.Title
// 			} else {
// 				// type is movie just need the title
// 				title = metadata.MediaContainer.Metadata[0].Title
// 			}

// 			user.Title = title
// 		}

// 		if err := db.Save(&user); err != nil {
// 			resp.Err = fmt.Sprintf("failed to save user: %v\n", err)
// 			resp.Write(w)
// 			return
// 		}

// 		resp.Result = user

// 		resp.Write(w)
// 	}
// }

// GetAllUsers returns all monitored users
// func GetAllUsers(db *storm.DB) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		var users []restrictedUser
// 		var resp clientResponse

// 		if err := db.All(&users); err != nil {
// 			resp.Err = fmt.Sprintf("failed to retrieve users: %v\n", err)
// 			resp.Write(w)
// 			return
// 		}

// 		resp.Result = users
// 		resp.Write(w)
// 	}
// }

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

// // GetPlexFriends will return an array of usernames and ids that are friends with associated plex token
// func GetPlexFriends(plexConnection *plex.Plex) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var resp clientResponse

// 		friends, err := plexConnection.GetFriends()

// 		if err != nil {
// 			resp.Err = fmt.Sprintf("failed to fetch friends from plex: %v", err)
// 			resp.Write(w)
// 			return
// 		}

// 		var friendsFiltered []plexFriend

// 		for _, friend := range friends {
// 			filteredFriend := plexFriend{
// 				ID:              strconv.Itoa(friend.ID),
// 				Username:        friend.Username,
// 				ServerID:        friend.Server.ID,
// 				ServerMachineID: friend.Server.MachineIdentifier,
// 				ServerName:      friend.Server.Name,
// 			}

// 			friendsFiltered = append(friendsFiltered, filteredFriend)
// 		}

// 		resp.Result = friendsFiltered
// 		resp.Write(w)
// 	}
// }

// // GetMetadataFromPlex fetches metadata of media from plex
// func GetMetadataFromPlex(plexConnection *plex.Plex) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var resp clientResponse

// 		mediaID := r.URL.Query().Get("mediaid")

// 		if mediaID == "" {
// 			resp.Err = "missing id from query: 'mediaID'"
// 			w.WriteHeader(http.StatusBadRequest)
// 			resp.Write(w)
// 			return
// 		}

// 		metadata, err := plexConnection.GetMetadataChildren(mediaID)

// 		if err != nil {
// 			resp.Err = fmt.Sprintf("failed to grab metadata from plex: %v", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			resp.Write(w)
// 			return
// 		}

// 		var results []plexSearchResults

// 		for _, child := range metadata.MediaContainer.Metadata {
// 			newResult := plexSearchResults{
// 				Title:     child.ParentTitle + " - " + child.Title,
// 				MediaID:   child.RatingKey,
// 				MediaType: child.Type,
// 			}

// 			results = append(results, newResult)
// 		}

// 		resp.Result = results

// 		resp.Write(w)
// 	}
// }

func cleanup(ctrlC chan os.Signal, shutdown chan bool, db store) {
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
func ConfigurePlexServer(db store, plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
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

				if err := db.savePlexServer(server{
					Name: "default",
					URL:  plexURL,
				}); err != nil {
					fmt.Printf("saving plex url failed: %v", err)
				}
			}
			if plexToken := r.FormValue("token"); plexToken != "" {
				plexConnection.setToken(plexToken)

				if err := db.savePlexToken(plexToken); err != nil {
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

	// connect to datastore
	db, err := initDataStore(storeDirectory)

	// handle interrupts - i.e. control-c
	go cleanup(ctrlC, shutdown, db)

	if err != nil {
		fmt.Printf("failed to open or create the datastore: %v\n", err)
		shutdown <- true
	}

	// check if app secret exists in datastore
	appSecret := []byte("iAmAseCReTuSEdTOENcrYp")

	if secret := db.getSecret(); len(secret) != 0 {
		db.secret = secret
	} else {
		// if not, create new appsecret and save
		// append a random number to make appsecret unique
		appSecret = append(appSecret, []byte(strconv.FormatInt(time.Now().Unix(), 10))...)

		if err := db.saveSecret(appSecret); err != nil {
			fmt.Printf("failed to save app secret: %v\n", err)
			shutdown <- true
		}

		db.secret = appSecret
	}

	// init plex stuff
	plexServerInfo, err := db.getPlexServer()

	if isVerbose && err != nil {
		fmt.Println("plex server:", err.Error())
	}

	if err != nil && err == badger.ErrKeyNotFound || err == badger.ErrEmptyKey {
		// init plex server in datastore
		if err := db.savePlexServer(server{}); err != nil {
			fmt.Printf("failed initializing plex server: %v", err)
		}
	} else if err != nil {
		fmt.Printf("failed retrieving plex server from data store: %v\n", err)
	}

	plexToken, err := db.getPlexToken()

	if isVerbose && err != nil {
		fmt.Println("plex token error:", err.Error())
	}

	if err != nil && err == badger.ErrKeyNotFound || err == badger.ErrEmptyKey {
		// init plex token in data store
		if err := db.savePlexServer(server{}); err != nil {
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

	// wh := plex.NewWebhook()

	// wh.OnPlay(OnPlay(db, plexConnection))

	// wh.OnStop(OnStop(db, plexConnection))

	// router.HandleFunc("/webhook", wh.Handler)

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/plex/server", ConfigurePlexServer(db, plexConnection)).Methods("POST", "PUT", "DELETE")

	// // add new restricted user
	// apiRouter.HandleFunc("/users/add", AddUser(db, plexConnection)).Methods("POST")

	// // list restricted users
	// apiRouter.HandleFunc("/users", GetAllUsers(db)).Methods("GET")

	// search media on plex
	apiRouter.HandleFunc("/search", SearchPlex(plexConnection)).Methods("GET")

	// get plex friends
	// apiRouter.HandleFunc("/friends", GetPlexFriends(plexConnection)).Methods("GET")

	// // get child data from plex
	// apiRouter.HandleFunc("/metadata", GetMetadataFromPlex(plexConnection)).Methods("GET")

	// fmt.Printf("serving one time plex on %s\n", config.Host)

	fmt.Println("server listening on " + *listenOnAddress)

	if err := http.ListenAndServe(*listenOnAddress, router); err != nil {
		fmt.Printf("server failed to start: %v\n", err)
		shutdown <- true
		return
	}
}
