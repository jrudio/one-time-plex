package datastore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

var isVerbose bool

// Store A wrapper around the database
type Store struct {
	db       *badger.DB
	isClosed bool
	keys     StoreKeys
	Secret   []byte
}

// StoreKeys keys for the database
type StoreKeys struct {
	appSecret  []byte
	plexToken  []byte
	plexServer []byte
	userPrefix []byte
	allUsers   []byte
}

// User a plex user that is attached to a media
type User struct {
	PlexUserID    string `json:"plexUserID"`
	Name          string `json:"plexUsername"`
	AssignedMedia `json:"assignedMedia"`
}

// AssignedMedia assigned media info such as: watch status, title, an key (plex media id)
type AssignedMedia struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// Serialize converts user into json
func (u User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// Server plex server info
type Server struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Serialize converts server info into json
func (s Server) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

// UnserializeServer converts json into Server
func UnserializeServer(serializedServer []byte) (Server, error) {
	var s Server

	err := json.Unmarshal(serializedServer, &s)

	return s, err
}

// UnserializeUser converts json into User
func UnserializeUser(user []byte) (User, error) {
	var u User

	err := json.Unmarshal(user, &u)

	return u, err
}

// InitDataStore creates a datastore in the specified directory
func InitDataStore(dirName string, _isVerbose bool) (Store, error) {
	var db Store

	if isVerbose {
		isVerbose = _isVerbose
		fmt.Println("checking if our datastore exists in the home directory at:", dirName)
	}

	// create a directory for our database
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if isVerbose {
			fmt.Println("creating directory because it doesn't exist")
		}

		if err := os.Mkdir(dirName, os.ModePerm); err != nil {
			return db, err
		}
	} else if !os.IsNotExist(err) && isVerbose {
		fmt.Println("datastore exists")
	}

	options := badger.DefaultOptions

	options.Dir = dirName
	options.ValueDir = dirName

	kvStore, err := badger.Open(options)

	if err != nil {
		return db, err
	}

	if isVerbose {
		fmt.Println("successfully opened data store")
	}

	db.db = kvStore
	db.keys = StoreKeys{
		appSecret:  []byte("app-secret"),
		plexToken:  []byte("plex-token"),
		plexServer: []byte("plex-server"),
		userPrefix: []byte("user-"), // holds the user info
		allUsers:   []byte("users"), // contains all user keys
	}

	return db, nil
}

// Close closes the datastore
func (s Store) Close() {
	if s.isClosed {
		fmt.Println("datastore already closed")
		return
	}

	if err := s.db.Close(); err != nil {
		fmt.Printf("datastore failed to closed: %v\n", err)
	}

	if isVerbose {
		fmt.Println("datastore is closed")
	}

	s.isClosed = true
}

// GetSecret fetches app secret
func (s Store) GetSecret() []byte {
	var secret []byte

	// an error is returned when the key is not found
	// so just return an empty secret
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.appSecret)

		if err != nil {
			return err
		}

		_secret, err := item.Value()

		if err != nil {
			return err
		}

		secret = _secret

		return nil
	})

	return secret
}

// SaveSecret saves the app secret
func (s Store) SaveSecret(secret []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.appSecret, secret)
	})
}

// GetPlexToken fetch and decrypt plex token
func (s Store) GetPlexToken() (string, error) {
	var plexToken string

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.plexToken)

		if err != nil {
			return err
		}

		tokenHash, err := item.Value()

		if err != nil {
			return err
		}

		_plexToken, err := decrypt(s.Secret, string(tokenHash))

		if err != nil {
			if isVerbose {
				fmt.Println("token decryption failed")
			}
			return err
		}

		plexToken = _plexToken

		return nil
	}); err != nil {
		return plexToken, err
	}

	if isVerbose {
		fmt.Printf("Your plex token is %s\n", plexToken)
	}

	return plexToken, nil
}

// SavePlexToken encrypt and save plex token in datastore
func (s Store) SavePlexToken(token string) error {
	tokenHash, err := encrypt(s.Secret, token)

	if err != nil {
		return err
	}

	if isVerbose {
		fmt.Printf("your plex token hash: %s\n", string(tokenHash))
	}

	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.plexToken, []byte(tokenHash))
	}); err != nil {
		return err
	}

	if isVerbose {
		fmt.Println("saved token hash to store")
	}

	return nil
}

// GetPlexServer fetches a plex server stored in the datastore
func (s Store) GetPlexServer() (Server, error) {
	var plexServer Server

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.plexServer)

		if err != nil {
			return err
		}

		serializedServer, err := item.Value()

		if err != nil {
			return err
		}

		_plexServer, err := UnserializeServer(serializedServer)

		if err != nil {
			return err
		}

		plexServer = _plexServer

		return nil
	})

	return plexServer, err
}

// SavePlexServer saves plex server info in the datastore
func (s Store) SavePlexServer(plexServer Server) error {
	serializedServer, err := plexServer.Serialize()
	if err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.plexServer, serializedServer)
	})
}

// SaveUser saves a user
func (s Store) SaveUser(user User) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := append(s.keys.userPrefix, []byte(user.PlexUserID)...)

		serializedUser, err := user.Serialize()

		if err != nil {
			return err
		}

		return txn.Set(key, serializedUser)
	})
}

// SaveUsers saves multiple users
func (s Store) SaveUsers(users []User) error {
	return s.db.Update(func(txn *badger.Txn) error {
		for _, user := range users {
			key := append(s.keys.userPrefix, []byte(user.PlexUserID)...)

			fmt.Println("saveusers key:", string(key))

			serializedUser, err := user.Serialize()

			if err != nil {
				return err
			}

			if err := txn.Set(key, serializedUser); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetUser fetches a user via id
func (s Store) GetUser(id string) (User, error) {
	var user User

	if err := s.db.View(func(txn *badger.Txn) error {
		key := append(s.keys.userPrefix, []byte(id)...)

		item, err := txn.Get(key)

		if err != nil {
			return err
		}

		serializedUser, err := item.Value()

		if err != nil {
			return err
		}

		_user, err := UnserializeUser(serializedUser)

		if err != nil {
			return err
		}

		user = _user

		return nil

	}); err != nil {
		return user, err
	}

	return user, nil
}

// GetAllUsers fetches all plex users that are assigned to media
func (s Store) GetAllUsers() (map[string]User, error) {
	users := map[string]User{}

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)

		defer it.Close()

		// prefix := append(s.keys.userPrefix, []byte("~")...)
		prefix := s.keys.userPrefix

		// fmt.Println("prefix:", string(prefix))

		for it.Seek(prefix); it.Valid(); it.Next() {
			item := it.Item()

			fmt.Println(string(item.Key()))

			serializedUser, err := item.Value()

			if err != nil {
				return err
			}

			user, err := UnserializeUser(serializedUser)

			if err != nil {
				return err
			}

			users[user.PlexUserID] = user
		}

		return nil
	})

	return users, err
}

// DeleteUser removes a user from the datastore
func (s Store) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		key := append(s.keys.userPrefix, []byte(id)...)

		return txn.Delete(key)
	})

	return err
}

// DeleteUsers removes multiple users from the datastore
func (s Store) DeleteUsers(userIDs []string) error {
	idLen := len(userIDs)

	err := s.db.Update(func(txn *badger.Txn) error {
		for i := 0; i < idLen; i++ {
			key := append(s.keys.userPrefix, []byte(userIDs[i])...)

			if err := txn.Delete(key); err != nil {
				fmt.Printf("failed to delete user id %s: %v\n", userIDs[i], err)
				continue
			}
		}

		return nil
	})

	return err
}
