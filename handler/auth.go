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

type TestRequest struct {
	Authenticated bool `json:"authenticated"`
}

func (a *Api) Test(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.String(500, "")
	} else {
		return c.JSON(200, TestRequest{Authenticated: context.IsAuthenticated()})
	}
}

func (a *Api) Logout(c echo.Context) error {
	if context, err := middleware.GetSessionContext(c); err != nil {
		return c.String(500, "")
	} else {
		if context.IsAuthenticated() {
			if err := middleware.Logout(a.DB, c); err != nil {
				if errors.Is(err, echo.ErrCookieNotFound) {
					return c.String(200, "")
				} else {
					return c.JSON(500, ErrorReturn{Error: "Database Error"})
				}
			}
		}
		return c.String(200, "")
	}
}

type LoginRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

type ErrorReturn struct {
	Error string `json:"error"`
}

func (a *Api) Login(c echo.Context) error {
	var f LoginRequest
	if err := utility.ValidateJsonForm(c, &f); err != nil {
		return c.JSON(400, ErrorReturn{Error: err.Error()})
	}

	if user, err := auth.Authenticate(a.DB, *f.Username, *f.Password); err != nil {
		return c.String(401, "")
	} else {
		if err := middleware.Login(a.DB, user, c); err != nil {
			return c.JSON(500, ErrorReturn{Error: err.Error()})
		}
	}

	return c.String(200, "")
}

type RegisterRequest struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *Api) Register(c echo.Context) error {
	var f RegisterRequest
	if err := utility.ValidateJsonForm(c, &f); err != nil {
		return c.JSON(400, ErrorReturn{Error: err.Error()})
	}

	var userCount int64
	if err := a.DB.Find(&utilitymodels.User{}, "username = ?", *f.Username).Count(&userCount).Error; err != nil {
		c.Logger().Error(err)
		return c.JSON(500, ErrorReturn{Error: "Database error"})
	}

	if userCount != 0 {
		return c.JSON(409, ErrorReturn{Error: "User with that username already exists"})
	}

	if _, err := database.CreateUser(a.DB, *f.Username, *f.Password, nil, true); err != nil {
		return c.JSON(500, ErrorReturn{Error: err.Error()})
	} else {
		return c.String(201, "")
	}
}

func (a *API) Connect(c echo.Context) error {
	_ = c // TODO
	return nil
}
