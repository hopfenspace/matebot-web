package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/models"
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
	Username         string               `json:"username"`
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

type listResponse struct {
	Message string       `json:"message"`
	Users   []simpleUser `json:"users"`
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
			Username:         coreUser.Name,
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

func (a *API) ListUsers(c echo.Context) error {
	coreUser, _, err := a.getUsers(c)
	if err != nil {
		return nil
	}
	if coreUser.Privilege() < MateBotSDKGo.Vouched {
		return c.JSON(400, GenericResponse{Message: "You are not permitted to request all users."})
	}
	u, err := a.SDK.GetUsers(map[string]string{"active": "true", "community": "false"})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	users := make([]simpleUser, len(u))
	for i := range u {
		var b models.CoreUser
		a.DB.Find(&b, "core_id = ?", u[i].ID)
		if b.ID == 0 {
			users[i] = simpleUser{
				UserID:   nil,
				CoreID:   u[i].ID,
				Username: u[i].Name,
			}
		} else {
			users[i] = simpleUser{
				UserID:   &b.UserID,
				CoreID:   u[i].ID,
				Username: u[i].Name,
			}
		}
	}
	return c.JSON(200, listResponse{Message: "OK", Users: users})
}
