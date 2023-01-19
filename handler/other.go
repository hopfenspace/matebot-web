package handler

import (
	"fmt"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
)

type consumablesResponse struct {
	Message     string                     `json:"message"`
	Consumables []*MateBotSDKGo.Consumable `json:"consumables"`
}

func (a *API) Consumables(c echo.Context) error {
	if _, _, err := a.getUnverifiedCoreID(c); err != nil {
		return err
	}
	consumables, err := a.SDK.GetConsumables(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, consumablesResponse{Message: "OK", Consumables: consumables})
}

type applicationsResponse struct {
	Message      string         `json:"message"`
	Applications []*namedObject `json:"applications"`
}

func (a *API) Applications(c echo.Context) error {
	if _, _, err := a.getUnverifiedCoreID(c); err != nil {
		return err
	}
	applications, err := a.SDK.GetApplications(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	convertedApplications := make([]*namedObject, len(applications))
	for i, a := range applications {
		convertedApplications[i] = &namedObject{
			ID:   a.ID,
			Name: a.Name,
		}
	}
	return c.JSON(200, applicationsResponse{Message: "OK", Applications: convertedApplications})
}

type balanceResponse struct {
	Message          string  `json:"message"`
	UserID           *uint64 `json:"user_id"`
	Username         *string `json:"username"`
	Balance          int64   `json:"balance"`
	BalanceFormatted string  `json:"balance_formatted"`
}

func (a *API) Balance(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	if *r.ID == coreUser.ID {
		return c.JSON(200, balanceResponse{Message: "OK", UserID: r.ID, Username: &coreUser.Name, Balance: coreUser.Balance, BalanceFormatted: a.SDK.FormatBalance(coreUser.Balance)})
	}
	if coreUser.Privilege() < MateBotSDKGo.Vouched {
		return c.JSON(400, GenericResponse{Message: "You are not permitted to request another user's balance."})
	}
	user, err := a.SDK.GetUser(int(*r.ID), nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, balanceResponse{Message: "OK", UserID: &user.ID, Username: &user.Name, Balance: user.Balance, BalanceFormatted: a.SDK.FormatBalance(user.Balance)})
}

type blameResponse struct {
	Message          string  `json:"message"`
	NobodyAvailable  bool    `json:"nobody_available"`
	UserID           *uint64 `json:"user_id"`
	Username         *string `json:"username"`
	Balance          *int64  `json:"balance"`
	BalanceFormatted *string `json:"balance_formatted"`
}

func (a *API) Blame(c echo.Context) error {
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	sponsor, err := a.SDK.FindSponsoringUser(coreUser)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	} else if sponsor == nil {
		return c.JSON(200, blameResponse{
			Message:          "Good news, nobody has to be blamed",
			NobodyAvailable:  true,
			UserID:           nil,
			Username:         nil,
			Balance:          nil,
			BalanceFormatted: nil,
		})
	}
	bF := a.SDK.FormatBalance(sponsor.Balance)
	return c.JSON(200, blameResponse{
		Message:          "There is a user who shall reduce his debts",
		NobodyAvailable:  false,
		UserID:           &sponsor.ID,
		Username:         &sponsor.Name,
		Balance:          &sponsor.Balance,
		BalanceFormatted: &bF,
	})
}

func (a *API) Zwegat(c echo.Context) error {
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	balance, err := a.SDK.GetCommunityBalance(coreUser)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	var msg string
	if balance >= 0 {
		msg = fmt.Sprintf("Peter errechnet ein massives Vermögen von %s!", a.SDK.FormatBalance(balance))
	} else {
		msg = fmt.Sprintf("Peter errechnet Gesamtschulden von %s!", a.SDK.FormatBalance(-balance))
	}
	return c.JSON(200, balanceResponse{Message: msg, UserID: nil, Username: nil, Balance: balance, BalanceFormatted: a.SDK.FormatBalance(balance)})
}
