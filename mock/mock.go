package mock

import "github.com/jrudio/otp-reloaded"

// UserService mocks the user service interface to allow for testing
type UserService struct {
	UserInvoked      bool
	UserFn           func(id string) (*otp.User, error)
	UsersFn          func() ([]*otp.User, error)
	CreateUserFn     func(u *otp.User) error
	DeleteUserFn     func(id string) error
	ChangePasswordFn func(id, newPassword string) error
	ChangeUsernameFn func(id, newName string) error
}

// User mocks User
func (u *UserService) User(id string) (*otp.User, error) {
	u.UserInvoked = true
	return u.UserFn(id)
}

// Users mocks Users
func (u UserService) Users() ([]*otp.User, error) {
	u.UserInvoked = true
	return u.UsersFn()
}

// CreateUser mocks CreateUser
func (u UserService) CreateUser(usr *otp.User) error {
	u.UserInvoked = true
	return u.CreateUserFn(usr)
}

// DeleteUser mocks DeleteUser
func (u UserService) DeleteUser(id string) error {
	u.UserInvoked = true
	return u.DeleteUserFn(id)
}

// ChangePassword mocks ChangePassword
func (u UserService) ChangePassword(id, newPassword string) error {
	u.UserInvoked = true
	return u.ChangePasswordFn(id, newPassword)
}

// ChangeUsername mocks ChangeUsername
func (u UserService) ChangeUsername(id, newName string) error {
	u.UserInvoked = true
	return u.ChangeUsernameFn(id, newName)
}
