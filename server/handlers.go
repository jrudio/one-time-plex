package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	plex "github.com/jrudio/go-plex-client"
	"github.com/jrudio/one-time-plex/server/datastore"
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

type addUserPost struct {
	UserID  string `json:"plexUserID"`
	MediaID string `json:"mediaID"`
}

func (a addUserPost) unserialize(data []byte) error {
	return json.Unmarshal(data, &a)
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
	w.Header().Set("Access-Control-Allow-Headers", "POST GET PUT")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")

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

		var postData addUserPost

		if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
			resp.Err = "failed to parse form data"
			resp.Write(w, http.StatusBadRequest)
			return
		}

		// check for required parameters
		// if err := r.ParseForm(); err != nil {
		// 	resp.Err = "failed to parse form"
		// 	resp.Write(w, http.StatusBadRequest)
		// 	return
		// }

		// plexUserID := r.FormValue("plexUserID")
		// mediaID := r.FormValue("mediaID")

		user.PlexUserID = postData.UserID
		user.AssignedMedia.ID = postData.MediaID
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

		title := metadata.Title

		if metadata.Type == "episode" {
			// title of the show + episode name
			title = results.MediaContainer.Title1 + ": " + metadata.Title
		}

		filteredMetadata[i] = plexMetadataChildren{
			Title:     title,
			Year:      "N/A",
			MediaID:   metadata.RatingKey,
			MediaType: metadata.Type,
		}

		// if year := strconv.FormatInt(metadata.Year, 10); year != "" {
		// 	filteredMetadata.Year = year
		// }

	}

	return filteredMetadata
}

func redactPlexPinInfo(pin plex.PinResponse) plex.PinResponse {
	pin.ID = 0
	pin.AuthToken = ""
	pin.Location = plex.PinResponse{}.Location
	pin.ClientIdentifier = ""

	return pin
}

// GetPlexPin get a plex pin from plex.tv
func GetPlexPin(db datastore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := clientResponse{}

		// check if we are already have a plex token
		plexToken, err := db.GetPlexToken()

		if err == nil {
			// and check if the token is valid
			plexConn, err := plex.New("", plexToken)

			if err != nil {
				fmt.Printf("GetPlexPin() auth token found in datastore: %v\n", err)
				response.Err = fmt.Sprintf("found auth token in datastore: %v", err)
				response.Write(w, http.StatusBadRequest)
				return
			}
			// fmt.Printf("could not get plex token: %v\n", err)

			if _, err := plexConn.MyAccount(); err != nil {
				fmt.Printf("GetPlexPin() auth token found in datastore: %v\n", err)
				response.Err = fmt.Sprintf("found auth token in datastore: %v", err)
				response.Write(w, http.StatusBadRequest)
				return
			}

			// the token is authorized so do not request plex pin
			response.Err = "we are already authorized; not fetching plex pin"
			response.Write(w, http.StatusUnprocessableEntity)
			return
		}

		// check to see if we already have a stored (and valid) pin from plex.tv
		plexPin, err := db.GetPlexPin()

		if err == nil {
			// redact some info
			plexPin = redactPlexPinInfo(plexPin)

			response.Result = plexPin
			response.Write(w, http.StatusOK)
			return
		}

		// grab a pin from plex.tv
		pin, err := plex.RequestPIN()

		if err != nil {
			response.Err = fmt.Sprintf("failed to fetch plex pin: %v", err)
			response.Write(w, http.StatusInternalServerError)
			return
		}

		// store response
		if err := db.SavePlexPin(pin); err != nil {
			response.Err = fmt.Sprintf("failed to save plex pin response: %v", err)
			response.Write(w, http.StatusInternalServerError)
			return
		}

		response.Result = redactPlexPinInfo(pin)
		response.Write(w, http.StatusOK)
	}
}

func isPinExpired(pin plex.PinResponse) bool {
	// convert to utc as the plex timestamp has no offset
	now := time.Now().UTC()

	// fmt.Println(pin.ExpiresAt)

	parsedTimestamp, err := time.Parse(time.RFC3339, pin.ExpiresAt)

	if err != nil {
		fmt.Printf("isPinExpired() parse timestamp failed: %v\n", err)
		return false
	}

	fmt.Printf("expires at: %dh %dm; currently: %dh %dm\n", parsedTimestamp.Hour(), parsedTimestamp.Minute(), now.Hour(), now.Minute())

	return now.After(parsedTimestamp)
}

// CheckPlexPin gets plex pin from datastore (if any) and checks plex.tv if valid
// returns error if expired or not authorized yet
func CheckPlexPin(db datastore.Store, plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := clientResponse{}

		// before we check for a plex pin
		// check if we are already have a plex token because
		// we might be authorized
		plexToken, err := db.GetPlexToken()

		if err == nil {
			// and check if the token is valid
			plexConn, err := plex.New("", plexToken)

			if err != nil {
				fmt.Printf("CheckPlexPin() auth token found in datastore: %v\n", err)
				response.Err = fmt.Sprintf("found auth token in datastore: %v", err)
				response.Write(w, http.StatusBadRequest)
				return
			}
			// fmt.Printf("could not get plex token: %v\n", err)

			if _, err := plexConn.MyAccount(); err != nil {
				fmt.Printf("CheckPlexPin() auth token found in datastore: %v\n", err)
				response.Err = fmt.Sprintf("found auth token in datastore: %v", err)
				response.Write(w, http.StatusBadRequest)
				return
			}

			// the token is authorized so do not request plex pin
			response.Err = "we are already authorized; not fetching plex pin"
			response.Write(w, http.StatusUnprocessableEntity)
			return
		}

		// get plex pin info from db
		plexPin, err := db.GetPlexPin()

		// we don't have a pin in our store to check
		if err != nil {
			response.Err = fmt.Sprintf("could not find plex pin: %v", err)
			response.Write(w, http.StatusNotFound)
			return
		}

		// pin has expired
		if isPinExpired(plexPin) {
			// fmt.Println("the plex pin is expired; request a new pin")
			if err := db.ClearPlexPin(); err != nil {
				fmt.Printf("checkPin() failed to clear plex pin from store: %v\n", err)
			}
			response.Err = fmt.Sprintf("pin has expired")
			response.Write(w, http.StatusNotFound)
			return
		}

		// check if valid and send error to client if expired
		pin, err := plex.CheckPIN(plexPin.ID, plexConnection.ClientIdentifier)

		// either invalid/expired or still not linked
		if err != nil {
			msg := "unauthorized"

			// our isPinExpired() should catch an expired pin
			// so this should be an error for not yet or an unauthorized pin

			// plex error code 1020 == "Code not found or expired"
			if len(pin.Errors) > 0 && pin.Errors[0].Code == 1020 {
				msg = pin.Errors[0].Message
				fmt.Printf("checkPlexPin() pin (%s): %s\n", plexPin.Code, msg)

			}

			response.Err = msg
			response.Write(w, http.StatusUnauthorized)
			return
		}

		// if pin is authorized save plex token to datastore
		if err := db.SavePlexToken(pin.AuthToken); err != nil {
			response.Err = fmt.Sprintf("failed to save plex token")
			response.Write(w, http.StatusInternalServerError)
			return
		}

		response.Result = true
		response.Write(w, http.StatusOK)
	}
}

// GetPlexServers returns a list of plex servers linked to our plex auth token
// gives client error if we are not authorized
func GetPlexServers(db datastore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := clientResponse{}

		plexToken, err := db.GetPlexToken()

		if err != nil {
			response.Err = fmt.Sprintf("not authorized; no plex token found: %v", err)
			response.Write(w, http.StatusUnauthorized)
			return
		}

		plexConnection, err := plex.New("", plexToken)

		if err != nil {
			response.Err = fmt.Sprintf("create plex client instance failed: %v", err)
			response.Write(w, http.StatusUnauthorized)
			return
		}

		servers, err := plexConnection.GetServers()

		if err != nil {
			response.Err = fmt.Sprintf("could not get plex servers: %v", err)
			response.Write(w, http.StatusInternalServerError)
			return
		}

		serverCount := len(servers)

		filteredServers := make([]plex.PMSDevices, serverCount)

		for i := 0; i < serverCount; i++ {
			server := servers[i]

			accessToken := "xxxxxxxx"
			accessToken = accessToken + server.AccessToken[len(server.AccessToken)-3:]

			filteredServers[i] = plex.PMSDevices{
				Name:        server.Name,
				Connection:  server.Connection,
				AccessToken: accessToken,
			}
		}

		response.Result = filteredServers
		response.Write(w, http.StatusOK)
	}
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

// TestPlexConnection sends a request to the specified plex server
// checking or auth token
func TestPlexConnection(w http.ResponseWriter, r *http.Request) {
	var resp clientResponse

	query := r.URL.Query()

	plexURL := query.Get("url")
	plexToken := query.Get("token")

	if plexURL == "" || plexToken == "" {
		resp.Err = "url and token are required"
		resp.Write(w, http.StatusBadRequest)
		return
	}

	plexConnection, err := plex.New(plexURL, plexToken)

	if err != nil {
		resp.Err = "testing connection to plex failed: " + err.Error()
		resp.Write(w, http.StatusBadRequest)
		return
	}

	result, err := plexConnection.Test()

	if err != nil {
		resp.Err = err.Error()

		resp.Write(w, http.StatusInternalServerError)
		return
	}

	resp.Result = result
	resp.Write(w, http.StatusOK)
}

type plexServerResp struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	URL   string `json:"url"`
}

// ConfigurePlexServer retrieve, add, edit, and delete plex server from one time plex
func ConfigurePlexServer(db datastore.Store, plexConnection *plexServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var response clientResponse

		switch r.Method {
		case "GET":
			server, err := db.GetPlexServer()

			if err != nil {
				response.Err = "failed to get plex server url: " + err.Error()
				response.Write(w, http.StatusInternalServerError)
				return
			}

			token, err := db.GetPlexToken()

			if err != nil {
				response.Err = "failed to get plex server token: " + err.Error()
				response.Write(w, http.StatusInternalServerError)
				return
			}

			serverResp := plexServerResp{
				URL:   server.URL,
				Token: token,
			}

			if err != nil {
				response.Err = "failed to serialize plex server: " + err.Error()
				response.Write(w, http.StatusInternalServerError)
				return
			}

			response.Result = serverResp
		case "POST":
			if err := r.ParseForm(); err != nil {
				response.Err = fmt.Sprintf("parsing body failed: %v", err)
				response.Write(w, http.StatusBadRequest)
				return
			}

			plexURL := r.PostFormValue("url")
			plexToken := r.PostFormValue("token")

			response.Result = fmt.Sprintf("received %s, %s", plexURL, plexToken)
		case "PUT":
			plexToken, err := db.GetPlexToken()

			if err != nil {
				response.Err = fmt.Sprintf("not authorized; no plex token found: %v", err)
				response.Write(w, http.StatusUnauthorized)
				return
			}

			plexTestConnection, err := plex.New("", plexToken)

			if err != nil {
				response.Err = fmt.Sprintf("create plex client instance failed: %v", err)
				response.Write(w, http.StatusUnauthorized)
				return
			}

			servers, err := plexTestConnection.GetServers()

			if err != nil {
				response.Err = fmt.Sprintf("could not get plex servers: %v", err)
				response.Write(w, http.StatusInternalServerError)
				return
			}

			var serverInfo plexServerResp

			if r.Header.Get("Accept") == "application/json" {

				if err := json.NewDecoder(r.Body).Decode(&serverInfo); err != nil {
					response.Err = "body parse failed: " + err.Error()
					response.Write(w, http.StatusBadRequest)
					return
				}
			} else {
				if err := r.ParseForm(); err != nil {
					response.Err = fmt.Sprintf("parsing body failed: %v", err)
					response.Write(w, http.StatusBadRequest)
					return
				}

				serverInfo.Name = r.FormValue("name")
				serverInfo.URL = r.FormValue("url")
			}

			// match requested server via name and url
			requestMatchesKnownServer := false

			for _, server := range servers {
				if requestMatchesKnownServer {
					break
				}

				if serverInfo.Name == server.Name {
					for _, connection := range server.Connection {
						if serverInfo.URL == connection.URI {
							requestMatchesKnownServer = true
							break
						}
					}
				}
			}

			if !requestMatchesKnownServer {
				response.Err = fmt.Sprintf("requested unknown server - url: %s name: %s", serverInfo.URL, serverInfo.Name)
				response.Write(w, http.StatusBadRequest)
				return
			}

			// set the plex connection with requested settings
			if serverInfo.URL != "" {
				plexConnection.setURL(serverInfo.URL)

				if err := db.SavePlexServer(datastore.Server{
					Name: serverInfo.Name,
					URL:  serverInfo.URL,
				}); err != nil {
					fmt.Printf("saving plex url failed: %v", err)
				}
			} else {
				response.Err = "url is required"
				response.Write(w, http.StatusBadRequest)
				return
			}

			if plexToken != "" {
				plexConnection.setToken(plexToken)

				// if err := db.SavePlexToken(plexToken); err != nil {
				// 	fmt.Printf("saving plex token failed: %v", err)
				// }
			} else {
				response.Err = "token is required"
				response.Write(w, http.StatusBadRequest)
				return
			}

			fmt.Printf("successfully set plex server to token: %s - url %s\n", plexToken, serverInfo.URL)
			response.Result = true
		case "OPTIONS":
			response.Result = true
		default:
			response.Err = "unknown method"
			response.Write(w, http.StatusMethodNotAllowed)
			return
		}

		response.Write(w, http.StatusOK)
	}
}

func getFileExtension(filename string) string {
	ext := ""
	filenameLen := len(filename)

	// search for the last period -- which is the start of the file's extension
	for i := filenameLen; i > 0; i-- {
		curr := filename[i-1 : i]
		ext = curr + ext

		if curr == "." {
			break
		}
	}

	return ext
}

func detectMimeType(filename string) string {
	switch getFileExtension(filename) {
	case ".json":
		return "application/json"
	case ".ico":
		return "image/x-icon"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".svg":
		return "image/svg+xml"
	default:
		return "text/plain"
	}
}

// ServeAsset uses go-bindata-assetfs' assets to serve
func ServeAsset(assetPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := Asset(assetPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(string("failed to access asset " + assetPath + ": " + err.Error())))
			return
		}

		contentType := detectMimeType(assetPath)

		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(a)
	}
}

// OnPlaying receives info when the 'play' event from the plex server
// it verifies with our database this user is allowed to play that media
func OnPlaying(react chan plexUserNotification, plexConnection *plexServer) func(n plex.NotificationContainer) {
	return func(n plex.NotificationContainer) {
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

		durationParsed, err := strconv.Atoi(metadata.MediaContainer.Metadata[0].Duration)

		if err != nil {
			fmt.Printf("duration parse failed: %v\n", err)
			return
		}

		duration = int64(durationParsed)
		title := metadata.MediaContainer.Metadata[0].Title

		currentTimeToSeconds := time.Duration(currentTime) * time.Millisecond
		durationToSeconds := time.Duration(duration) * time.Millisecond

		sessions, err := plexConnection.GetSessions()

		if err != nil {
			fmt.Printf("failed to fetch sessions on plex server: %v\n", err)
			return
		}

		for _, session := range sessions.MediaContainer.Metadata {
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
	}
}
