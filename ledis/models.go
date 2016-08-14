package ledis

import (
	"errors"
	"fmt"
	"github.com/jrudio/otp-reloaded"
	"github.com/rs/xid"
	"github.com/siddontang/ledisdb/ledis"
	"golang.org/x/crypto/bcrypt"
)

const userKey = "user"
const userlistKey = "userlist"

// users are stored in a hash map:
// key: "user:abc123"
// 	field: "username"
// 	value: "bob"
// 	field: "password"
// 	value: "@#$4t"
//
// A list of users are stored in set as a Go map that has been JSON-serialized
// key: "userlist"
// '{"<id>": "<empty-value>"}'
//
// This helps Users() retrieve all users in a key-value datastore

// UserService manages user actions
type UserService struct {
	DB *ledis.DB
}

// User returns a user object of the requested id
func (u UserService) User(id string) (*otp.User, error) {
	userFields, err := u.DB.HMget([]byte(userKey+":"+id), []byte("username"))

	if err != nil {
		return &otp.User{}, err
	}

	var usr otp.User

	usr.ID = id
	usr.Name = string(userFields[0])

	return &usr, nil
}

// Users returns all users
func (u UserService) Users() ([]*otp.User, error) {
	userIDList, err := u.getUserlist()

	if err != nil {
		return []*otp.User{}, err
	}

	users := make([]*otp.User, len(userIDList))

	listIndex := 0

	for username, userID := range userIDList {
		usr := &otp.User{
			ID:   userID,
			Name: username,
		}

		users[listIndex] = usr

		listIndex++
	}

	return users, nil
}

// CreateUser creates a new user
func (u UserService) CreateUser(usr *otp.User) (string, error) {
	nameExists, err := u.usernameExists(usr.Name)

	if err != nil {
		return "", err
	}

	if nameExists {
		return "", errors.New("username exists")
	}

	usr.ID = xid.New().String()

	var passwordHash []byte
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(usr.Password), 12)

	if err != nil {
		return "", err
	}

	// save user info
	err = u.DB.HMset([]byte(userKey+":"+usr.ID), ledis.FVPair{
		Field: []byte("username"),
		Value: []byte(usr.Name),
	},
		ledis.FVPair{
			Field: []byte("password"),
			Value: passwordHash,
		},
		ledis.FVPair{
			Field: []byte("apiKey"),
			Value: []byte(usr.APIKey),
		})

	if err != nil {
		return "", err
	}

	// add user id to userlist
	if err = u.addIDToUserlist(usr.ID, usr.Name); err != nil {
		// remove user if we couldn't add to list

		_, err2 := u.DB.HClear([]byte(userKey + ":" + usr.ID))

		return "", fmt.Errorf("addIDToUserList: %v\nHClear: %v", err, err2)
	}

	return usr.ID, nil
}

// DeleteUser removes a user from the database
func (u UserService) DeleteUser(id string) error {
	return nil
}

// ChangePassword changes a user's password
func (u UserService) ChangePassword(id, newPassword string) error {
	return nil
}

// ChangeUsername changes a user's username
func (u UserService) ChangeUsername(id, newName string) error {
	return nil
}

func (u UserService) usernameExists(username string) (bool, error) {
	id, err := u.DB.HGet([]byte(userlistKey), []byte(username))

	if err != nil {
		return false, err
	}

	return string(id) != "", nil
}

func (u UserService) getUserlist() (map[string]string, error) {
	userlist, err := u.DB.HGetAll([]byte(userlistKey))

	if err != nil {
		return map[string]string{}, err
	}

	userlistMap := map[string]string{}

	for _, user := range userlist {
		// username: userid
		userlistMap[string(user.Field)] = string(user.Value)
	}

	return userlistMap, err
}

func (u UserService) addIDToUserlist(id, username string) error {
	_, err := u.DB.HSet([]byte(userlistKey), []byte(username), []byte(id))

	return err
}
