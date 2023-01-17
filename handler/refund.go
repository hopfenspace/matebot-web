package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
	"strconv"
)

type refundResponse struct {
	Message string  `json:"message"`
	Refund  *refund `json:"refund"`
}

type refundsResponse struct {
	Message string   `json:"message"`
	Refunds []refund `json:"refunds"`
}

type refundVoteResponse struct {
	Message string  `json:"message"`
	Vote    *vote   `json:"vote"`
	Refund  *refund `json:"refund"`
}

type refund struct {
	ID              uint         `json:"id"`
	Amount          uint         `json:"amount"`
	AmountFormatted string       `json:"amount_formatted"`
	Description     string       `json:"description"`
	CreatorID       uint         `json:"creator_id"`
	CreatorName     string       `json:"creator_name"`
	Active          bool         `json:"active"`
	Created         uint         `json:"created"`
	Modified        uint         `json:"modified"`
	Allowed         *bool        `json:"allowed"`
	BallotID        uint         `json:"ballot_id"`
	Votes           []vote       `json:"votes"`
	Transaction     *transaction `json:"transaction"`
}

func (a *API) convRefund(r *MateBotSDKGo.Refund) *refund {
	var votes []vote
	for i := range r.Votes {
		votes[i] = *a.convVote(r.Votes[i])
	}
	return &refund{
		ID:              r.ID,
		Amount:          r.Amount,
		AmountFormatted: a.SDK.FormatBalance(int(r.Amount)),
		Description:     r.Description,
		CreatorID:       r.Creator.ID,
		CreatorName:     r.Creator.Name,
		Active:          r.Active,
		Created:         *r.Created,
		Modified:        *r.Modified,
		Allowed:         r.Allowed,
		BallotID:        r.BallotID,
		Votes:           votes,
		Transaction:     a.convTransaction(r.Transaction),
	}
}

func (a *API) NewRefund(c echo.Context) error {
	var r newMoneyRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	refund, err := a.SDK.NewRefund(coreUserID, *r.Amount, *r.Description)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, refundResponse{Message: "OK", Refund: a.convRefund(refund)})
}

func (a *API) voteOnRefund(c echo.Context, approve bool) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	refunds, err := a.SDK.GetRefunds(map[string]string{"id": strconv.Itoa(int(*r.ID))})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	response, err := a.SDK.VoteOnRefundBallot(refunds[0].BallotID, coreUserID, approve)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, refundVoteResponse{Message: "OK", Refund: a.convRefund(&response.Refund), Vote: a.convVote(response.Vote)})
}

func (a *API) ApproveRefund(c echo.Context) error {
	return a.voteOnRefund(c, true)
}

func (a *API) DisapproveRefund(c echo.Context) error {
	return a.voteOnRefund(c, false)
}

func (a *API) AbortRefund(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	refund, err := a.SDK.AbortRefund(*r.ID, coreUserID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, refundResponse{Message: "OK", Refund: a.convRefund(refund)})
}

func (a *API) OpenRefunds(c echo.Context) error {
	r, err := a.SDK.GetRefunds(map[string]string{"active": "true"})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	refunds := make([]refund, len(r))
	for i := range r {
		refunds[i] = *a.convRefund(r[i])
	}
	return c.JSON(200, refundsResponse{Message: "OK", Refunds: refunds})
}

func (a *API) AllRefunds(c echo.Context) error {
	coreUser, _, err := a.getUnverifiedCoreUser(c)
	if err != nil {
		return nil
	}
	if coreUser.Privilege() < MateBotSDKGo.Internal {
		return c.JSON(400, GenericResponse{Message: "You are not permitted to request all refunds."})
	}
	r, err := a.SDK.GetRefunds(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	refunds := make([]refund, len(r))
	for i := range r {
		refunds[i] = *a.convRefund(r[i])
	}
	return c.JSON(200, refundsResponse{Message: "OK", Refunds: refunds})
}
