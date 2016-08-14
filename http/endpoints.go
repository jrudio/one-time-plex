package http

import "github.com/labstack/echo"

// CreateEndpoints starts the api server
func CreateEndpoints(h Handler) *echo.Echo {
	router := echo.New()

	api := router.Group("/api/v1")

	// user management
	// this will likely not stay in the build as
	// I want to keep this lean. By lean I mean there
	//  will only be one user and that user will be the admin
	//
	// The reason for having a user at all is for the api key
	//
	//
	// user := api.Group("/user")

	// user.GET("/:id", h.getUser)
	// user.POST("/new", h.createUser)

	// api.GET("/users", h.getUsers)

	// login
	// user.POST("/login")
	// logout / invalidate token
	// user.GET("/logout")
	// user.PUT("/password")
	// user.PUT("/name")

	// protected routes
	// request := api.Group("/request", checkAPIKey)
	request := api.Group("/request")

	request.POST("/add", h.addUserToPMS)

	return router
}
