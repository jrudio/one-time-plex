package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

type dbWrapper struct {
	db      Store
	tempDir string
}

func (d dbWrapper) Close() {
	defer os.RemoveAll(d.tempDir)

	d.db.Close()
}

// createDatastore temporarily creates a datastore for testing
func createDatastore() (dbWrapper, error) {
	var wrapper dbWrapper

	dir, err := ioutil.TempDir("", ".one-time-plex")

	if err != nil {
		return wrapper, err
	}

	wrapper.tempDir = dir

	db, err := InitDataStore(dir, false)

	if err != nil {
		return wrapper, err
	}

	wrapper.db = db

	return wrapper, nil
}

var users map[string]User

func TestPollUsers(t *testing.T) {
	db, err := createDatastore()

	if err != nil {
		t.Errorf("failed creating temp datastore: %v", err)
	}

	// seed database
	seedUserCount := 30
	seedUsers := make([]User, seedUserCount)

	fmt.Println("seed length:", seedUserCount)

	for i := 0; i < seedUserCount; i++ {
		seedUsers[i] = User{
			Name:       "test",
			PlexUserID: "test-" + strconv.Itoa(i),
			AssignedMedia: AssignedMedia{
				ID:     "abc123",
				Title:  "fake title",
				Status: "not started",
			},
		}
	}

	if err := db.db.SaveUsers(seedUsers); err != nil {
		t.Errorf("failed to seed users: %v", err)
		return
	}

	defer db.Close()

	_users, err := db.db.GetAllUsers()

	if err != nil {
		t.Errorf("failed retrieving users: %v", err)
	}

	for i := 0; i < seedUserCount; i++ {
		id := "test-" + strconv.Itoa(i)

		if _, ok := _users[id]; !ok {
			t.Errorf("couldn't find id: %s", id)
		}
	}
}

func BenchmarkPollUsers(b *testing.B) {
	db, err := createDatastore()

	if err != nil {
		b.Errorf("failed creating temp datastore: %v", err)
	}

	// seed database
	seedUsers := make([]User, 30)

	for i := 0; i < 30; i++ {
		seedUsers[i] = User{
			Name:       "test",
			PlexUserID: "bench-" + strconv.Itoa(i),
			AssignedMedia: AssignedMedia{
				ID:    "abc123",
				Title: "fake title",
			},
		}
	}

	if err := db.db.SaveUsers(seedUsers); err != nil {
		b.Errorf("failed to seed users: %v", err)
		return
	}

	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_users, err := db.db.GetAllUsers()

		if err != nil {
			b.Errorf("failed retrieving users: %v", err)
		}

		users = _users
	}
}
