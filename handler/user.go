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
	UserID           uint                 `json:"user_id"`
	CoreID           uint                 `json:"core_id"`
	Balance          int                  `json:"balance"`
	BalanceFormatted string               `json:"balance_formatted"`
	Permission       bool                 `json:"permission"`
	Active           bool                 `json:"active"`
	External         bool                 `json:"external"`
	VoucherId        interface{}          `json:"voucher_id"`
	Aliases          []MateBotSDKGo.Alias `json:"aliases"`
	Created          uint                 `json:"created"`
	Modified         uint                 `json:"modified"`
}

type stateResponse struct {
	User    user   `json:"user"`
	Message string `json:"message"`
}

func (a *API) Me(c echo.Context) error {
	coreUser, localUser, err := a.getUsers(c)
	if err != nil {
		return err
	}
	return c.JSON(200, stateResponse{
		Message: "OK",
		User: user{
			UserID:           localUser.ID,
			CoreID:           coreUser.ID,
			Balance:          coreUser.Balance,
			BalanceFormatted: a.SDK.FormatBalance(coreUser.Balance),
			Permission:       coreUser.Permission,
			Active:           coreUser.Active,
			External:         coreUser.External,
			VoucherId:        coreUser.VoucherID,
			Aliases:          coreUser.Aliases,
			Created:          coreUser.Created,
			Modified:         coreUser.Modified,
		},
	})
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
