package mock

import "github.com/jrudio/go-plex-client"

// PlexMonitorService mocks the user service interface to allow for testing
type PlexMonitorService struct {
	UserInvoked bool
	UserFn      func(id string) (plex.MonitoredUser, error)
	UsersFn     func() ([]plex.MonitoredUser, error)
	AddUserFn   func(id int, ratingKey string) error
}

// User mocks User
func (u *PlexMonitorService) User(id string) (plex.MonitoredUser, error) {
	u.UserInvoked = true
	return u.UserFn(id)
}

// Users mocks Users
func (u PlexMonitorService) Users() ([]plex.MonitoredUser, error) {
	u.UserInvoked = true
	return u.UsersFn()
}

// AddUser mocks AddUser
func (u PlexMonitorService) AddUser(id int, ratingKey string) error {
	u.UserInvoked = true
	return u.AddUserFn(id, ratingKey)
}
