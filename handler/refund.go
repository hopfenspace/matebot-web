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
	ID              uint64       `json:"id"`
	Amount          uint64       `json:"amount"`
	AmountFormatted string       `json:"amount_formatted"`
	Description     string       `json:"description"`
	CreatorID       uint64       `json:"creator_id"`
	CreatorName     string       `json:"creator_name"`
	Active          bool         `json:"active"`
	Created         uint64       `json:"created"`
	Modified        uint64       `json:"modified"`
	Allowed         *bool        `json:"allowed"`
	BallotID        uint64       `json:"ballot_id"`
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
		AmountFormatted: a.SDK.FormatBalance(int64(r.Amount)),
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
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	refund, err := a.SDK.NewRefund(coreUser.ID, *r.Amount, *r.Description)
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
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	refunds, err := a.SDK.GetRefunds(map[string]string{"id": strconv.Itoa(int(*r.ID))})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	response, err := a.SDK.VoteOnRefundBallot(refunds[0].BallotID, coreUser.ID, approve)
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
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	refund, err := a.SDK.AbortRefund(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, refundResponse{Message: "OK", Refund: a.convRefund(refund)})
}

func (a *API) OpenRefunds(c echo.Context) error {
	l := MateBotSDKGo.Vouched
	_, _, err := a.getVerifiedCoreUser(c, &l)
	if err != nil {
		return nil
	}
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
	l := MateBotSDKGo.Internal
	_, _, err := a.getVerifiedCoreUser(c, &l)
	if err != nil {
		return nil
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
