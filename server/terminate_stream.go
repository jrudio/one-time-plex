package main

import (
	"fmt"

	"github.com/jrudio/one-time-plex/server/datastore"
)

func streamGuard(react <-chan plexUserNotification, db datastore.Store, plexConnection *plexServer, hasPlexPass bool) {
	if isVerbose {
		fmt.Println("firing up plex event listener...")
	}

	for {
		select {
		case r := <-react:
			fmt.Printf("user id: %s, rating key (media id): %s, and current time: %d\n", r.userID, r.ratingKey, r.currentTime)
			userInStore, err := db.GetUser(r.userID)

			if err != nil {
				fmt.Printf("failed to retrieve user from datastore with id of %s\n", r.userID)
				continue
			}

			// end stream
			if userInStore.AssignedMedia.ID != r.ratingKey {
				fmt.Printf("\t%s (%s) is not allowed to watch %s\n\tattempting to terminate stream...\n", userInStore.Name, userInStore.PlexUserID, r.ratingKey)

				if hasPlexPass {
					// plex pass - will stop playback on whatever player the user is using
					if err := plexConnection.TerminateSession(r.sessionID, "You are not allowed to watch this"); err != nil {
						fmt.Printf("failed to terminate session id %s: %v\n", r.sessionID, err)
					}
				} else {
					fmt.Println("\tclient id:", r.clientID)

					// we cannot use StopPlayback if the user is watching via web browser
					// so we must remove their access to our library
					if r.playerType == "Plex Web" {
						fmt.Println("\tuser is watching via web player -- we must remove their access to library")
						// if err := removeFriendAccess(r.userID, plexConnection); err != nil {
						// 	fmt.Printf("failed to revoke friend's access (%s):  %v", r.userID, err)
						// 	continue
						// }

						// userInStore.IsFriend = false

						// if err := db.SaveUser(userInStore); err != nil {
						// 	fmt.Printf("saving user failed when trying to modify IsFriend on %s: %v\n",
						// 		userInStore.Name,
						// 		err)
						// }
						continue
					}

					// non-plex pass - we can only prevent the user from downloading more data from server
					if plexConnection.StopPlayback(r.clientID); err != nil {
						fmt.Printf("failed to stop playback for user %s (%s), session %s\n", userInStore.Name, r.userID, r.sessionID)
					}
				}

			}
		}
	}

}

// func removeFriendAccess(userid string, plexConnection *plexServer) error {
// 	_, err := plexConnection.RemoveFriendAccessToLibrary(userid, plexConnection.machineID, plexConnection.serverID)

// 	return err
// }
