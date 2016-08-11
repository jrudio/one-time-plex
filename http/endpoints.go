package http

import (
	"github.com/labstack/echo"
)

// CreateEndpoints starts the api server
func CreateEndpoints(h Handler) *echo.Echo {
	router := echo.New()

	api := router.Group("/api/v1")

	// user management
	user := api.Group("/user")

	user.GET("/:id", h.getUser)

	// login
	// user.POST("/login")
	// logout / invalidate token
	// user.GET("/logout")
	// user.PUT("/password")
	// user.PUT("/name")

	// protected routes

	return router
}
