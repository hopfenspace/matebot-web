package handler

import (
	"fmt"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
	"strconv"
)

type pollResponse struct {
	Message string `json:"message"`
	Poll    *poll  `json:"poll"`
}

type pollsResponse struct {
	Message string `json:"message"`
	Polls   []poll `json:"polls"`
}

type pollVoteResponse struct {
	Message string `json:"message"`
	Vote    *vote  `json:"vote"`
	Poll    *poll  `json:"poll"`
}

type poll struct {
	ID          uint                     `json:"id"`
	Active      bool                     `json:"active"`
	Accepted    *bool                    `json:"accepted"`
	Variant     MateBotSDKGo.PollVariant `json:"variant"`
	UserID      uint                     `json:"user_id"`
	UserName    string                   `json:"user_name"`
	CreatorID   uint                     `json:"creator_id"`
	CreatorName string                   `json:"creator_name"`
	Created     uint                     `json:"created"`
	Modified    uint                     `json:"modified"`
	BallotID    uint                     `json:"ballot_id"`
	Votes       []vote                   `json:"votes"`
}

func (a *API) convPoll(p *MateBotSDKGo.Poll) *poll {
	votes := make([]vote, len(p.Votes))
	for i := range p.Votes {
		votes[i] = *a.convVote(p.Votes[i])
	}
	users, err := a.SDK.GetUsers(map[string]string{"id": strconv.Itoa(int(p.CreatorID))})
	if err != nil {
		return nil
	}
	return &poll{
		ID:          p.ID,
		Active:      p.Active,
		Accepted:    p.Accepted,
		Variant:     p.Variant,
		UserID:      p.User.ID,
		UserName:    p.User.Name,
		CreatorID:   p.CreatorID,
		CreatorName: users[0].Name,
		Created:     p.Created,
		Modified:    p.Modified,
		BallotID:    p.BallotID,
		Votes:       votes,
	}
}

type newPoll struct {
	User    *any                      `json:"user" echotools:"required"`
	Variant *MateBotSDKGo.PollVariant `json:"variant" echotools:"required;not empty"`
}

func (a *API) NewPoll(c echo.Context) error {
	var r newPoll
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	variant, found := map[MateBotSDKGo.PollVariant]MateBotSDKGo.PollVariant{
		MateBotSDKGo.GetInternal:     MateBotSDKGo.GetInternal,
		MateBotSDKGo.LooseInternal:   MateBotSDKGo.LooseInternal,
		MateBotSDKGo.GetPermission:   MateBotSDKGo.GetPermission,
		MateBotSDKGo.LoosePermission: MateBotSDKGo.LoosePermission,
	}[*r.Variant]
	fmt.Println("r", r, *r.User, *r.Variant, found)
	if !found {
		return c.JSON(400, GenericResponse{Message: fmt.Sprintf("Invalid poll variant '%s'", *r.Variant)})
	}
	switch (*r.User).(type) {
	case string:
		poll, err := a.SDK.NewPoll((*r.User).(string), coreUserID, string(variant))
		if err != nil {
			return c.JSON(400, GenericResponse{Message: err.Error()})
		}
		return c.JSON(200, pollResponse{Message: "OK", Poll: a.convPoll(poll)})
	case float64:
		poll, err := a.SDK.NewPoll(int((*r.User).(float64)), coreUserID, string(variant))
		if err != nil {
			return c.JSON(400, GenericResponse{Message: err.Error()})
		}
		return c.JSON(200, pollResponse{Message: "OK", Poll: a.convPoll(poll)})
	default:
		return c.JSON(400, GenericResponse{Message: "Unknown JSON format for user"})
	}
}

func (a *API) voteOnPoll(c echo.Context, approve bool) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	polls, err := a.SDK.GetPolls(map[string]string{"id": strconv.Itoa(int(*r.ID))})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	response, err := a.SDK.VoteOnPollBallot(polls[0].BallotID, coreUserID, approve)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, pollVoteResponse{Message: "OK", Poll: a.convPoll(&response.Poll), Vote: a.convVote(response.Vote)})
}

func (a *API) ApprovePoll(c echo.Context) error {
	return a.voteOnPoll(c, true)
}

func (a *API) DisapprovePoll(c echo.Context) error {
	return a.voteOnPoll(c, false)
}

func (a *API) AbortPoll(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUserID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil
	}
	poll, err := a.SDK.AbortPoll(*r.ID, coreUserID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, pollResponse{Message: "OK", Poll: a.convPoll(poll)})
}

func (a *API) OpenPolls(c echo.Context) error {
	p, err := a.SDK.GetPolls(map[string]string{"active": "true"})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	polls := make([]poll, len(p))
	for i := range p {
		polls[i] = *a.convPoll(p[i])
	}
	return c.JSON(200, pollsResponse{Message: "OK", Polls: polls})
}

func (a *API) AllPolls(c echo.Context) error {
	coreUser, _, err := a.getUnverifiedCoreUser(c)
	if err != nil {
		return nil
	}
	if coreUser.Privilege() < MateBotSDKGo.Internal {
		return c.JSON(400, GenericResponse{Message: "You are not permitted to request all polls."})
	}
	p, err := a.SDK.GetPolls(nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	polls := make([]poll, len(p))
	for i := range p {
		polls[i] = *a.convPoll(p[i])
	}
	return c.JSON(200, pollsResponse{Message: "OK", Polls: polls})
}
