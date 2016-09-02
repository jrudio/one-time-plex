package otp

import "github.com/jrudio/go-plex-client"

// otp glues together the packages and sub-packages

// the domain of this app entails:
//  - Handling user access to Plex
//  - Act upon user behavior

// User is a user that is authorized to talk to our server
type User struct {
	ID       string `json:"id" form:"id"`
	Name     string `json:"username" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
	APIKey   string `json:"apiKey,omitempty" form:"apiKey"`
}

// MediaRequest holds the information required to request access to media
type MediaRequest struct {
	ID string `json:"id"`
	// RatingKey is the id of the media used by Plex
	RatingKey string `json:"ratingKey"`
	PlexUser
}

// PlexUser is a plex account which will be used to invite a user to PMS
type PlexUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

// UserService interfaces with users
type UserService interface {
	User(id string) (*User, error)
	UserViaName(username string) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) (string, error)
	DeleteUser(id string) error
	ChangePassword(id, newPassword string) error
	ChangeUsername(id, newName string) error
}

// MediaRequestService interfaces with media requests
type MediaRequestService interface {
	Request(id string) (*MediaRequest, error)
	Requests() ([]*MediaRequest, error)
	CreateRequest(m *MediaRequest) error
	DeleteRequest(id string) error
	StopRequestInProgress(id string) error
}

// PlexMonitorService interfaces with the monitor package of jrudio/go-plex-client
type PlexMonitorService interface {
	Userlist() (map[int]plex.MonitoredUser, error)
	User(id int) (plex.MonitoredUser, error)
	UserViaName(username string) (plex.MonitoredUser, error)
	AddUser(id int, username, ratingKey string) error
	SetField(id int, field, value string) error
	Start() error
	Stop() error
}
