package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
)

type communismResponse struct {
	Message   string     `json:"message"`
	Communism *communism `json:"communism"`
}

type communismsResponse struct {
	Message    string      `json:"message"`
	Communisms []communism `json:"communisms"`
}

type communism struct {
	ID               uint64                              `json:"id"`
	Amount           uint64                              `json:"amount"`
	Description      string                              `json:"description"`
	CreatorID        uint64                              `json:"creator_id"`
	CreatorName      string                              `json:"creator_name"`
	Active           bool                                `json:"active"`
	Created          uint64                              `json:"created"`
	Modified         uint64                              `json:"modified"`
	Participants     []MateBotSDKGo.CommunismParticipant `json:"participants"`
	MultiTransaction *multiTransaction                   `json:"multi_transaction"`
}

func (a *API) convCommunism(c *MateBotSDKGo.Communism) *communism {
	var mT *multiTransaction
	if c.MultiTransaction != nil {
		transactions := make([]transaction, len(c.MultiTransaction.Transactions))
		for i, t := range c.MultiTransaction.Transactions {
			transactions[i] = *a.convTransaction(&t)
		}
		transactions[0] = transaction{}
		mT = &multiTransaction{
			BaseAmount:   c.MultiTransaction.BaseAmount,
			TotalAmount:  c.MultiTransaction.TotalAmount,
			Transactions: transactions,
			Timestamp:    c.MultiTransaction.Timestamp,
		}
	}
	return &communism{
		ID:               c.ID,
		Amount:           c.Amount,
		Description:      c.Description,
		CreatorID:        c.CreatorID,
		Active:           c.Active,
		Created:          c.Created,
		Modified:         c.Modified,
		Participants:     c.Participants,
		MultiTransaction: mT,
	}
}

func (a *API) NewCommunism(c echo.Context) error {
	var r newMoneyRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	communism, err := a.SDK.NewCommunism(coreUser.ID, *r.Amount, *r.Description)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, communismResponse{Message: "OK", Communism: a.convCommunism(communism)})
}

func (a *API) CloseCommunism(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	communism, err := a.SDK.CloseCommunism(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, communismResponse{Message: "OK", Communism: a.convCommunism(communism)})
}

func (a *API) JoinCommunism(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	communism, err := a.SDK.IncreaseCommunismParticipation(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, communismResponse{Message: "OK", Communism: a.convCommunism(communism)})
}

func (a *API) LeaveCommunism(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	communism, err := a.SDK.DecreaseCommunismParticipation(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, communismResponse{Message: "OK", Communism: a.convCommunism(communism)})
}

func (a *API) AbortCommunism(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	communism, err := a.SDK.AbortCommunism(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, communismResponse{Message: "OK", Communism: a.convCommunism(communism)})
}

func (a *API) OpenCommunisms(c echo.Context) error {
	l := MateBotSDKGo.Vouched
	_, _, err := a.getVerifiedCoreUser(c, &l)
	if err != nil {
		return nil
	}
	co, err := a.SDK.GetCommunisms(map[string]string{"active": "true"})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	var communisms []communism
	for i := range co {
		communisms[i] = *a.convCommunism(co[i])
	}
	return c.JSON(200, communismsResponse{Message: "OK", Communisms: communisms})
}

func (a *API) AllCommunisms(c echo.Context) error {
	l := MateBotSDKGo.Internal
	_, _, err := a.getVerifiedCoreUser(c, &l)
	if err != nil {
		return nil
	}
	co, err := a.SDK.GetCommunisms(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	var communisms []communism
	for i := range co {
		communisms[i] = *a.convCommunism(co[i])
	}
	return c.JSON(200, communismsResponse{Message: "OK", Communisms: communisms})
}
