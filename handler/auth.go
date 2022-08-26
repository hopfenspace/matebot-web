package handler

import (
	"errors"
	"fmt"
	"github.com/hopfenspace/matebot-web/models"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
)

type TestResponse struct {
	Message       string `json:"message"`
	Authenticated bool   `json:"authenticated" json:"authenticated,omitempty"`
}

func (a *API) Test(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.JSON(500, TestResponse{Authenticated: false, Message: "Session error"})
	} else {
		if context.IsAuthenticated() {
			return c.JSON(200, TestResponse{Authenticated: true, Message: "Successfully authenticated"})
		} else {
			return c.JSON(200, TestResponse{Authenticated: false, Message: "Not authenticated"})
		}
	}
}

func (a *API) Logout(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.JSON(500, GenericResponse{Message: "Invalid session"})
	} else {
		if context.IsAuthenticated() {
			if err := middleware.Logout(a.DB, c); err != nil {
				if errors.Is(err, echo.ErrCookieNotFound) {
					return c.JSON(200, GenericResponse{Message: "Successfully logged out"})
				} else {
					return c.JSON(500, GenericResponse{Message: "Database Error"})
				}
			}
		}
		return c.JSON(200, GenericResponse{Message: "Successfully logged out"})
	}
}

type LoginRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *API) Login(c echo.Context) error {
	var r LoginRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}

	if user, err := auth.AuthenticateLocalUser(a.DB, *r.Username, *r.Password); err != nil {
		return c.JSON(401, GenericResponse{Message: "Invalid username or password"})
	} else {
		if err := middleware.Login(a.DB, user, c, true); err != nil {
			return c.JSON(500, GenericResponse{Message: err.Error()})
		}
	}
	return c.JSON(200, GenericResponse{Message: "Successfully logged in"})
}

type RegisterRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *API) Register(c echo.Context) error {
	var r RegisterRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}

	var userCount int64
	if err := a.DB.Find(&utilitymodels.LocalUser{}, "username = ?", *r.Username).Count(&userCount).Error; err != nil {
		c.Logger().Error(err)
		return c.JSON(500, GenericResponse{Message: "Database error"})
	}
	if userCount != 0 {
		return c.JSON(409, GenericResponse{Message: "User with that username already exists"})
	}

	coreUser, err := a.SDK.NewUserWithAlias(*r.Username)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}

	localUser, err := database.CreateLocalUser(a.DB, *r.Username, *r.Password, nil)
	if err != nil {
		return c.JSON(500, GenericResponse{Message: err.Error()})
	}

	u := models.CoreUser{
		UserID:    localUser.ID,
		MateBotID: coreUser.ID,
	}
	if err := a.DB.Create(&u).Error; err != nil {
		return c.JSON(500, GenericResponse{Message: err.Error()})
	}

	c.Logger().Infof("Registered new user %s (core %d, local %d)", *r.Username, coreUser.ID, localUser.ID)
	return c.JSON(201, GenericResponse{Message: "Successfully registered new account"})
}

type ConnectRequest struct {
	Username         *string `json:"username" echotools:"required;not empty"`
	Password         *string `json:"password" echotools:"required;not empty"`
	ExistingUsername *string `json:"existing_username" echotools:"required;not empty"`
	Application      *string `json:"application" echotools:"required;not empty"`
}

func (a *API) Connect(c echo.Context) error {
	var r ConnectRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}

	var userCount int64
	if err := a.DB.Find(&utilitymodels.LocalUser{}, "username = ?", *r.Username).Count(&userCount).Error; err != nil {
		c.Logger().Error(err)
		return c.JSON(500, GenericResponse{Message: "Database error"})
	}
	if userCount != 0 {
		return c.JSON(409, GenericResponse{Message: "User with that username already exists"})
	}

	users, err := a.SDK.GetUsers(map[string]string{"active": "true", "alias_username": *r.ExistingUsername, "alias_application": *r.Application})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	} else if len(users) == 0 {
		return c.JSON(400, GenericResponse{Message: fmt.Sprintf("No user '%s' found", *r.ExistingUsername)})
	} else if len(users) > 1 {
		return c.JSON(409, GenericResponse{Message: fmt.Sprintf("Multiple users '%s' found", *r.ExistingUsername)})
	}

	user := users[0]
	alias, err := a.SDK.NewAlias(user.ID, *r.Username)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}

	localUser, err := database.CreateLocalUser(a.DB, *r.Username, *r.Password, nil)
	if err != nil {
		return c.JSON(500, GenericResponse{Message: err.Error()})
	}

	u := models.CoreUser{
		UserID:    localUser.ID,
		MateBotID: user.ID,
	}
	if err := a.DB.Create(&u).Error; err != nil {
		return c.JSON(500, GenericResponse{Message: err.Error()})
	}

	c.Logger().Infof("Connected user %s (core %d, local %d) with new alias ID %d", *r.Username, user.ID, localUser.ID, alias.ID)
	return c.JSON(201, GenericResponse{Message: "Successfully registered new account"})
}
