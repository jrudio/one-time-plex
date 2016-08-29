package http

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/go-plex-client"
	"github.com/jrudio/otp-reloaded"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

// Handler holds the endpoint handlers
type Handler struct {
	UserService        otp.UserService
	MediaService       otp.MediaRequestService
	PlexMonitorService otp.PlexMonitorService
	Plex               *plex.Plex
}

type endpointResponse struct {
	Err    bool        `json:"error"`
	Reason string      `json:"reason"`
	Result interface{} `json:"result"`
}

// Serve starts a web server with endpoints
func (h Handler) ServeHTTP(router *echo.Echo, listenAddress string) {
	router.Run(standard.New(listenAddress))
}

func (h Handler) getUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.UserService.User(id)

	if err != nil {
		log.WithField("endpoint", c.Request().URI()).Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if user.ID == "" {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (h Handler) getUsers(c echo.Context) error {
	users, err := h.UserService.Users()

	if err != nil {
		log.WithField("endpoint", c.Request().URI()).Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, endpointResponse{
		Result: users,
	})
}

func (h Handler) createUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing username")
	}

	if password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing password")
	}

	usr := otp.User{
		Name:     username,
		Password: password,
	}

	userID, err := h.UserService.CreateUser(&usr)

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": c.Request().URI(),
			"username": username,
			"password": password,
		}).Error(err)

		return echo.NewHTTPError(http.StatusExpectationFailed, "failed to create user: "+err.Error())
	}

	return c.JSON(http.StatusOK, endpointResponse{
		Result: userID,
	})
}

func (h Handler) addUserToPMS(c echo.Context) error {
	resp := endpointResponse{
		Err: true,
	}

	ratingKey := c.FormValue("ratingKey")

	if ratingKey == "" {
		resp.Reason = "missing ratingKey"
		return c.JSON(http.StatusBadRequest, resp)
	}

	plexPass := c.QueryParam("plexpass")
	guest := c.QueryParam("guest")

	if plexPass == "1" {
		// use labels
	}

	var plexUsername string

	if guest == "1" {
		// supply a guest Plex account
	} else {
		plexUsername = c.FormValue("plexUsername")

		if plexUsername == "" {
			resp.Reason = "missing plex username"
			return c.JSON(http.StatusBadRequest, resp)
		}

		// validate the username via Plex's validate endpoint
		// isValid, err := h.Plex.CheckUsernameOrEmail(plexUsername)

		// if err != nil {
		// 	log.WithFields(log.Fields{
		// 		"endpoint":     c.Request().URI(),
		// 		"ratingKey":    ratingKey,
		// 		"plexUsername": plexUsername,
		// 	}).Error(err)
		// 	resp.Reason = "could not check plex username"
		// 	return c.JSON(http.StatusBadRequest, resp)
		// }

		// if !isValid {
		// 	resp.Reason = "plex username is not valid"
		// 	return c.JSON(http.StatusBadRequest, resp)
		// }
	}

	// get library id of media via ratingKey
	mediaInfo, err := h.Plex.GetMetadata(ratingKey)

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
		}).Error(err)
		resp.Reason = "could not grab information for key " + ratingKey
		return c.JSON(http.StatusBadRequest, resp)
	}

	libraryKey := mediaInfo.LibrarySectionID

	var machineID string
	var libraryID int
	var serverSections []plex.ServerSections

	// get machine id
	machineID, err = h.Plex.GetMachineID()

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
		}).Error(err)

		resp.Reason = "could not grab the PMS machine id"
		return c.JSON(http.StatusBadRequest, resp)
	}

	serverSections, err = h.Plex.GetSections(machineID)

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
			"machine id:":  machineID,
		}).Error(err)
		resp.Reason = "could not get server sections"
		return c.JSON(http.StatusBadRequest, resp)
	}

	// get library id which is a required param to invite user
	for _, section := range serverSections {
		if libraryKey != section.Key {
			continue
		}

		libraryID = section.ID
	}

	if libraryID == 0 {
		resp.Reason = "could not get library id"
		return c.JSON(http.StatusBadRequest, resp)
	}

	var plexFriendUserID int

	// we need the id of the user to monitor
	// if we are inviting the user then that will be straightforward
	// if the user is on the server then we'll compare usernames and select the matching user id
	invite := c.QueryParam("invite")

	if invite == "1" {
		// invite user to pms
		inviteParams := plex.InviteFriendParams{
			LibraryIDs:      []int{libraryID},
			MachineID:       machineID,
			UsernameOrEmail: plexUsername,
		}

		if plexPass == "1" {
			inviteParams.Label = plexUsername
		}
		plexFriendUserID, err = h.Plex.InviteFriend(inviteParams)

		if err != nil {
			log.Info("inviteFriend")
			log.WithFields(log.Fields{
				"endpoint":     c.Request().URI(),
				"ratingKey":    ratingKey,
				"plexUsername": plexUsername,
				"machine id":   machineID,
			}).Error(err)
			resp.Reason = "failed invite user"
			return c.JSON(http.StatusBadRequest, resp)
		}
	} else {
		// search users for user id
		friends, err := h.Plex.GetFriends()

		if err != nil {
			log.Info("getFriends")
			log.WithFields(log.Fields{
				"endpoint":     c.Request().URI(),
				"ratingKey":    ratingKey,
				"plexUsername": plexUsername,
				"machine id":   machineID,
			}).Error(err)
			resp.Reason = "failed getting friends list"
			return c.JSON(http.StatusInternalServerError, resp)
		}

		for _, friend := range friends {
			if friend.Username == plexUsername {
				plexFriendUserID = friend.ID
			}
		}
	}

	if plexFriendUserID == 0 {
		log.Info("plexFriendUserID")
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
			"machine id":   machineID,
		}).Error(err)
		resp.Reason = fmt.Sprintf("failed getting id of %s", plexUsername)
		return c.JSON(http.StatusInternalServerError, resp)
	}

	// monitor user
	if err := h.PlexMonitorService.AddUser(plexFriendUserID, ratingKey); err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
			"machine id":   machineID,
		}).Error(err)

		resp.Reason = "failed to add user to monitor list"

		return c.JSON(http.StatusExpectationFailed, resp)
	}

	resp.Err = false
	resp.Result = true

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) startPlexMonitor(c echo.Context) error {
	resp := endpointResponse{
		Err: true,
	}

	if err := h.PlexMonitorService.Start(); err != nil {
		resp.Reason = "failed to start plex monitor: " + err.Error()
		c.JSON(http.StatusExpectationFailed, resp)
		return nil
	}

	resp.Err = false
	resp.Result = true

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) stopPlexMonitor(c echo.Context) error {
	resp := endpointResponse{
		Err: true,
	}

	if err := h.PlexMonitorService.Stop(); err != nil {
		resp.Reason = "failed to stop plex monitor: " + err.Error()
		return c.JSON(http.StatusExpectationFailed, resp)
	}

	resp.Err = false
	resp.Result = true

	return c.JSON(http.StatusOK, resp)
}
