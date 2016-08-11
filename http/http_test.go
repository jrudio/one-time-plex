package http

import (
	"github.com/jrudio/otp-reloaded"
	"github.com/jrudio/otp-reloaded/mock"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerUser(t *testing.T) {
	var userSrvc mock.UserService
	var h Handler
	h.UserService = &userSrvc

	userSrvc.UserFn = func(id string) (*otp.User, error) {
		if id != "abc123" {
			t.Fatalf("unexpected id: %s", id)
		}

		return &otp.User{
			ID:     "abc123",
			Name:   "bob",
			APIKey: "abc123",
		}, nil
	}

	e := echo.New()
	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
	c.SetPath("/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc123")

	err := h.getUser(c)

	if err != nil {
		t.Fatal(err)
	}

	if !userSrvc.UserInvoked {
		t.Fatal("expected User() to be invoked")
	}
}
