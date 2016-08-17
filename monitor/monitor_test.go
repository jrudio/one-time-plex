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
		t.Errorf("user id is %d instead of expected 2\n", user.UserID)
	}

	if user.RatingKey != "abc123" {
		t.Errorf("ratingkey is %s instead of expected 'abc123'\n", user.RatingKey)
	}
}
