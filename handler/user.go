package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
)

type simpleUser struct {
	UserID   *uint  `json:"user_id"`
	CoreID   uint   `json:"core_id"`
	Username string `json:"username"`
}

type user struct {
	UserID     uint                 `json:"user_id"`
	CoreID     uint                 `json:"core_id"`
	Balance    int                  `json:"balance"`
	Permission bool                 `json:"permission"`
	Active     bool                 `json:"active"`
	External   bool                 `json:"external"`
	VoucherId  interface{}          `json:"voucher_id"`
	Aliases    []MateBotSDKGo.Alias `json:"aliases"`
	Created    int                  `json:"created"`
	Modified   int                  `json:"modified"`
}

func (a *API) State(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) ChangeUsername(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) StartVouching(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) StopVouching(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) DropPrivileges(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) ConfirmAlias(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) DeleteAlias(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}
