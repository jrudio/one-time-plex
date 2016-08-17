package monitor

import (
	"fmt"
	"github.com/jrudio/otp-reloaded/ledis"
	ledisDB "github.com/siddontang/ledisdb/ledis"
	"os"
	"testing"
	// "github.com/pkg/errors"
)

var db *ledisDB.DB

func init() {
	var err error
	db, err = ledis.Open()

	if err != nil {
		fmt.Println("could not open database: ", err)
		os.Exit(1)
	}
}

func TestUserList(t *testing.T) {
	db.FlushAll()

	plexService := PlexMonitorService{
		DB: db,
	}

	// seed list
	for ii := 0; ii < 4; ii++ {
		if err := plexService.AddUser(ii, "abc"); err != nil {
			t.Error(err)
		}
	}

	users, errUserlist := plexService.Userlist()

	if errUserlist != nil {
		t.Fatalf("could not retrieve userlist: %v\n", errUserlist)
	}

	userCount := len(users)

	if userCount != 4 {
		t.Error("expected 4 users recieved: ", userCount)
	}
}

func TestUser(t *testing.T) {
	db.FlushAll()

	plexService := PlexMonitorService{
		DB: db,
	}

	if err := plexService.AddUser(2, "abc123"); err != nil {
		t.Fatalf("failed to add user: %v\n", err)
	}

	user, err := plexService.User(2)

	if err != nil {
		t.Errorf("failed to retrieve user: %v\n", err)
	}

	if user.UserID != 2 {
		t.Errorf("user id is '%d' instead of expected 2\n", user.UserID)
	}

	if user.RatingKey != "abc123" {
		t.Errorf("ratingkey is '%s' instead of expected 'abc123'\n", user.RatingKey)
	}
}

func TestSetUser(t *testing.T) {
	db.FlushAll()

	plexService := PlexMonitorService{
		DB: db,
	}

	if err := plexService.AddUser(5, "abc123"); err != nil {
		t.Errorf("failed to add user: %v\n", err)
	}

	if err := plexService.SetUserField(5, "isTranscoding", "1"); err != nil {
		t.Errorf("failed to set isTranscoding field: %v\n", err)
	}

	if err := plexService.SetUserField(5, "isDirectPlay", "1"); err != nil {
		t.Errorf("failed to set isTranscoding field: %v\n", err)
	}

	if err := plexService.SetUserField(5, "killingSession", "1"); err != nil {
		t.Errorf("failed to set isTranscoding field: %v\n", err)
	}

	if err := plexService.SetUserField(5, "killedSession", "1"); err != nil {
		t.Errorf("failed to set isTranscoding field: %v\n", err)
	}

	user, err2 := plexService.User(5)

	if err2 != nil {
		t.Errorf("failed to retrieve user: %v\n", err2)
		return
	}

	if !user.IsTranscode {
		t.Errorf("user id %d isTranscode should be 'true'\n", user.UserID)
	}

	if !user.IsDirectPlay {
		t.Errorf("user id %d isDirectPlay should be 'true'\n", user.UserID)
	}

	if !user.KilledSession {
		t.Errorf("user id %d killedSession should be 'true'\n", user.UserID)
	}

	if !user.KillingSession {
		t.Errorf("user id %d killingSession should be 'true'\n", user.UserID)
	}
}

func BenchmarkUser(b *testing.B) {
	db.FlushAll()

	plexService := PlexMonitorService{
		DB: db,
	}

	if err := plexService.AddUser(2, "abc123"); err != nil {
		fmt.Printf("failed to add user: %v\n", err)
	}

	b.ResetTimer()

	for ii := 0; ii < b.N; ii++ {
		_, err := plexService.User(2)

		if err != nil {
			fmt.Printf("failed to retrieve user: %v\n", err)
		}
	}
}
