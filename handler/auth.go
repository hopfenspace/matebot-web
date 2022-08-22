package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
)

type TestResponse struct {
	Error         bool   `json:"error"`
	Message       string `json:"message"`
	Authenticated bool   `json:"authenticated" json:"authenticated,omitempty"`
}

func (a *API) Test(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.JSON(500, TestResponse{Authenticated: false, Error: true, Message: "Session error"})
	} else {
		return c.JSON(200, TestResponse{Authenticated: context.IsAuthenticated(), Error: false, Message: "Successfully authenticated"})
	}
}

func (a *API) Logout(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.JSON(500, GenericResponse{Error: true, Message: "Invalid session"})
	} else {
		if context.IsAuthenticated() {
			if err := middleware.Logout(a.DB, c); err != nil {
				if errors.Is(err, echo.ErrCookieNotFound) {
					return c.JSON(200, GenericResponse{Error: false, Message: "Successfully logged out"})
				} else {
					return c.JSON(500, GenericResponse{Error: true, Message: "Database Error"})
				}
			}
		}
		return c.JSON(200, GenericResponse{Error: false, Message: "Successfully logged out"})
	}
}

type LoginRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *API) Login(c echo.Context) error {
	var r LoginRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Error: true, Message: err.Error()})
	}

	if user, err := auth.AuthenticateLocalUser(a.DB, *r.Username, *r.Password); err != nil {
		return c.JSON(401, GenericResponse{Error: true, Message: "Invalid username or password"})
	} else {
		if err := middleware.Login(a.DB, user, c, true); err != nil {
			return c.JSON(500, GenericResponse{Error: true, Message: err.Error()})
		}
	}
	return c.JSON(200, GenericResponse{Error: false, Message: "Successfully logged in"})
}

type RegisterRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *API) Register(c echo.Context) error {
	var r RegisterRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Error: true, Message: err.Error()})
	}

	var userCount int64
	if err := a.DB.Find(&utilitymodels.LocalUser{}, "username = ?", *r.Username).Count(&userCount).Error; err != nil {
		c.Logger().Error(err)
		return c.JSON(500, GenericResponse{Error: true, Message: "Database error"})
	}

	if userCount != 0 {
		return c.JSON(409, GenericResponse{Error: true, Message: "User with that username already exists"})
	}

	if _, err := database.CreateLocalUser(a.DB, *r.Username, *r.Password, nil); err != nil {
		return c.JSON(500, GenericResponse{Error: true, Message: err.Error()})
	} else {
		return c.JSON(201, GenericResponse{Error: false, Message: "Successfully registered new account"})
	}
}

type ConnectRequest struct {
	Username    *string `json:"username" echotools:"required;not empty"`
	Password    *string `json:"password" echotools:"required;not empty"`
	Alias       *string `json:"alias" echotools:"required;not empty"`
	Application *string `json:"application" echotools:"required;not empty"`
}

func (a *API) Connect(c echo.Context) error {
	return nil
}
