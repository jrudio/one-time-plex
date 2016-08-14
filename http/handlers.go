package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/go-plex-client"
	"github.com/jrudio/otp-reloaded"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

// Handler holds the endpoint handlers
type Handler struct {
	UserService  otp.UserService
	MediaService otp.MediaRequestService
	Plex         *plex.Plex
}

type endpointResponse struct {
	Err    string      `json:"error"`
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
	ratingKey := c.FormValue("ratingKey")

	if ratingKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing ratingKey")
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
			return echo.NewHTTPError(http.StatusBadRequest, "missing plex username")
		}

		// validate the username via Plex's validate endpoint
		isValid, err := h.Plex.CheckUsernameOrEmail(plexUsername)

		if err != nil {
			log.WithFields(log.Fields{
				"endpoint":     c.Request().URI(),
				"ratingKey":    ratingKey,
				"plexUsername": plexUsername,
			}).Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "could not check plex username")
		}

		if !isValid {
			return echo.NewHTTPError(http.StatusBadRequest, "plex username is not valid")
		}
	}

	// get library id of media via ratingKey
	mediaInfo, err := h.Plex.GetMetadata(ratingKey)

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
		}).Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "could not grab information for key "+ratingKey)
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
		return echo.NewHTTPError(http.StatusBadRequest, "could not grab the PMS machine id")
	}

	serverSections, err = h.Plex.GetSections(machineID)

	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
			"machine id:":  machineID,
		}).Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "could not get server sections")
	}

	// get library id which is a required param to invite user
	for _, section := range serverSections {
		if libraryKey != section.Key {
			continue
		}

		libraryID = section.ID
	}

	if libraryID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "could not get library id")
	}

	// invite user to pms
	inviteParams := plex.InviteFriendParams{
		LibraryIDs:      []int{libraryID},
		MachineID:       machineID,
		UsernameOrEmail: plexUsername,
	}

	if plexPass == "1" {
		inviteParams.Label = plexUsername
	}

	var isFriendInvited bool
	isFriendInvited, err = h.Plex.InviteFriend(inviteParams)

	if err != nil {
		log.Info("inviteFriend")
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
			"machine id":   machineID,
		}).Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed invite user")
	}

	// monitor user

	// return success

	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"endpoint":     c.Request().URI(),
	// 		"ratingKey":    ratingKey,
	// 		"plexUsername": plexUsername,
	// 	}).Error(err)

	// 	return echo.NewHTTPError(http.StatusExpectationFailed, "failed to add user to PMS: "+err.Error())
	// }

	if !isFriendInvited {
		log.WithFields(log.Fields{
			"endpoint":     c.Request().URI(),
			"ratingKey":    ratingKey,
			"plexUsername": plexUsername,
		}).Debug("invite failed")

		return echo.NewHTTPError(http.StatusExpectationFailed, "failed to add user to PMS")
	}

	return c.JSON(http.StatusOK, endpointResponse{
		Result: true,
	})
}
