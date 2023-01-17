package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
	"strconv"
)

type consumablesResponse struct {
	Message     string                     `json:"message"`
	Consumables []*MateBotSDKGo.Consumable `json:"consumables"`
}

func (a *API) Consumables(c echo.Context) error {
	consumables, err := a.SDK.GetConsumables(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, consumablesResponse{Message: "OK", Consumables: consumables})
}

type applicationsResponse struct {
	Message      string                      `json:"message"`
	Applications []*MateBotSDKGo.Application `json:"applications"`
}

func (a *API) Applications(c echo.Context) error {
	applications, err := a.SDK.GetApplications(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, applicationsResponse{Message: "OK", Applications: applications})
}

type balanceResponse struct {
	Message          string  `json:"message"`
	UserID           *uint   `json:"user_id"`
	Username         *string `json:"username"`
	Balance          int     `json:"balance"`
	BalanceFormatted string  `json:"balance_formatted"`
}

func (a *API) Balance(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getUnverifiedCoreUser(c)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	if *r.ID == coreUser.ID {
		return c.JSON(200, balanceResponse{Message: "OK", UserID: r.ID, Username: &coreUser.Name, Balance: coreUser.Balance, BalanceFormatted: a.SDK.FormatBalance(coreUser.Balance)})
	}
	if coreUser.Privilege() < MateBotSDKGo.Internal {
		return c.JSON(400, GenericResponse{Message: "You are not permitted to request another user's balance."})
	}
	users, err := a.SDK.GetUsers(map[string]string{"id": strconv.Itoa(int(*r.ID))})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	} else if len(users) != 1 {
		c.Logger().Error(users, "required one element as 'users' response")
		return c.JSON(400, GenericResponse{Message: "Invalid response"})
	}
	return c.JSON(200, balanceResponse{Message: "OK", UserID: &users[0].ID, Username: &users[0].Name, Balance: users[0].Balance, BalanceFormatted: a.SDK.FormatBalance(users[0].Balance)})
}

func (a *API) Blame(c echo.Context) error {
	coreUser, _, err := a.getUnverifiedCoreUser(c)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	sponsor, err := a.SDK.FindSponsoringUser(coreUser)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, balanceResponse{Message: "OK", UserID: &sponsor.ID, Username: &sponsor.Name, Balance: sponsor.Balance, BalanceFormatted: a.SDK.FormatBalance(sponsor.Balance)})
}

func (a *API) Zwegat(c echo.Context) error {
	coreUser, _, err := a.getUnverifiedCoreUser(c)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	balance, err := a.SDK.GetCommunityBalance(coreUser)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, balanceResponse{Message: "OK", UserID: nil, Username: nil, Balance: balance, BalanceFormatted: a.SDK.FormatBalance(balance)})
}
