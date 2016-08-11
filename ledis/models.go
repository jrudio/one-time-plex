package ledis

import (
	"github.com/jrudio/otp-reloaded"
	"github.com/siddontang/ledisdb/ledis"
)

// UserService manages user actions
type UserService struct {
	DB *ledis.DB
}

// User returns a user object of the requested id
func (u UserService) User(id string) (*otp.User, error) {
	return &otp.User{}, nil
}

// Users returns all users
func (u UserService) Users() ([]*otp.User, error) {
	return []*otp.User{}, nil
}

// CreateUser creates a new user
func (u UserService) CreateUser(usr *otp.User) error {
	return nil
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
