package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jrudio/otp-reloaded"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

// Handler holds the endpoint handlers
type Handler struct {
	UserService  otp.UserService
	MediaService otp.MediaRequestService
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
