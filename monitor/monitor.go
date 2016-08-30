package monitor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/go-plex-client"
	"github.com/pkg/errors"
	"github.com/siddontang/ledisdb/ledis"
	"strconv"
)

const (
	plexUserlistKey = "plex.userlist"
	plexUserKey     = "plexUser:"
)

// PlexMonitorService will interface with PMS using any datastore
type PlexMonitorService struct {
	DB      *ledis.DB
	Plex    *plex.Plex
	Monitor plex.Monitor
}

// Count returns a count of a monitored users
func (p PlexMonitorService) Count() int {
	userlistHashMap, err := p.DB.HGetAll([]byte("userlist"))

	if err != nil {
		log.WithField("list count", "failed to get list").Error(err)
		return 0
	}

	userlistCount := len(userlistHashMap)

	return userlistCount
}

// Userlist returns a map of monitored users
func (p *PlexMonitorService) Userlist() (map[int]plex.MonitoredUser, error) {
	userlistHashMap, err := p.DB.HGetAll([]byte(plexUserlistKey))

	if err != nil {
		return map[int]plex.MonitoredUser{}, errors.Wrap(err, "failed to retrieve userlist")
	}

	userlistCount := len(userlistHashMap)

	userlist := make(map[int]plex.MonitoredUser, userlistCount)

	for _, userlistField := range userlistHashMap {
		userID, err := strconv.Atoi(string(userlistField.Field))

		if err != nil {
			return userlist, errors.Wrap(err, "failed to convert user id from string to int")
		}

		var user plex.MonitoredUser

		user, err = p.User(userID)

		if err != nil {
			return userlist, errors.Wrap(err, "failed to retrieve user")
		}

		userlist[userID] = user
	}

	return userlist, nil
}

// User returns a list of monitored users
func (p PlexMonitorService) User(id int) (plex.MonitoredUser, error) {
	idStr := strconv.Itoa(id)

	var user plex.MonitoredUser

	// make sure user exists
	bytesWritten, err := p.DB.HKeyExists([]byte(plexUserKey + idStr))

	if err != nil {
		return user, errors.Wrap(err, "check existing user failed")
	}

	if bytesWritten < 1 {
		return user, errors.New("user doesn't exist")
	}

	// extract required fields
	var fields [][]byte
	fields, err = p.DB.HMget([]byte(plexUserKey+idStr), []byte("ratingKey"), []byte("isTranscoding"), []byte("isDirectPlay"), []byte("killingSession"), []byte("killedSession"))

	if err != nil {
		return user, errors.Wrap(err, "failed to retrieve user")
	}

	user.RatingKey = string(fields[0])
	user.IsTranscode = string(fields[1]) == "1"
	user.IsDirectPlay = string(fields[2]) == "1"
	user.KillingSession = string(fields[3]) == "1"
	user.KilledSession = string(fields[4]) == "1"
	user.UserID = id

	return user, nil
}

// AddUser adds a user to the monitoring list
func (p PlexMonitorService) AddUser(id int, ratingKey string) error {
	idStr := strconv.Itoa(id)

	// check for existing user
	exists, err := p.DB.HKeyExists([]byte(plexUserKey + idStr))

	if err != nil {
		return errors.Wrap(err, "check existing user failed")
	}

	if exists > 0 {
		return errors.New("user already exists")
	}

	err = p.DB.HMset([]byte(plexUserKey+idStr), ledis.FVPair{
		Field: []byte("ratingKey"),
		Value: []byte(ratingKey),
	},
		ledis.FVPair{
			Field: []byte("isTranscoding"),
			Value: []byte("0"),
		},
		ledis.FVPair{
			Field: []byte("isDirectPlay"),
			Value: []byte("0"),
		},
		ledis.FVPair{
			Field: []byte("killingSession"),
			Value: []byte("0"),
		},
		ledis.FVPair{
			Field: []byte("killedSession"),
			Value: []byte("0"),
		})

	if err != nil {
		p.DB.HClear([]byte(idStr))

		return errors.Wrap(err, "add user failed")
	}

	// add to userlist
	var bytesWritten int64
	bytesWritten, err = p.DB.HSet([]byte(plexUserlistKey), []byte(idStr), []byte(""))

	if err != nil {
		p.DB.HClear([]byte(idStr))

		return errors.Wrap(err, "add to userlist failed")
	}

	if bytesWritten < 1 {
		p.DB.HClear([]byte(idStr))

		return errors.New("add to userlist failed")
	}

	return nil
}

// RemoveUser removes a user from being monitored
func (p PlexMonitorService) RemoveUser(id int) error {
	return nil
}

// SetField sets a field on a user
func (p PlexMonitorService) SetField(id int, field string, value string) error {
	idStr := strconv.Itoa(id)

	// check for existing user
	exists, err := p.DB.HKeyExists([]byte(plexUserKey + idStr))

	if err != nil {
		return errors.Wrap(err, "check existing user failed")
	}

	if exists < 1 {
		return errors.New("user doesn't exist")
	}

	_, err = p.DB.HSet([]byte(plexUserKey+idStr), []byte(field), []byte(value))

	if err != nil {
		return errors.Wrap(err, "failed to set user field")
	}

	return nil
}

// Start starts monitoring plex users
func (p *PlexMonitorService) Start() error {
	if p.Plex.URL == "" {
		return errors.New("plex connection is required")
	}

	if p.Monitor.Init {
		return errors.New("monitor has already started")
	}

	p.Monitor = plex.Monitor{
		Interval: 1500,
		Userlist: p,
		PlexConn: p.Plex,
	}

	p.Monitor.Start()

	return nil
}

// Stop stops monitoring users
func (p PlexMonitorService) Stop() error {
	return nil
}
