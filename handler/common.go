package handler

import (
	"errors"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/models"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type eventNotification struct {
	Type MateBotSDKGo.EventType
	Data any
}

type eventWrapper struct {
	receivers     *[]uint64
	minPrivilege  *MateBotSDKGo.PrivilegeLevel
	confirmedOnly bool
	notification  eventNotification
}

type eventChannelKey struct {
	sessionID string
	confirmed bool
	coreID    uint64
	privilege MateBotSDKGo.PrivilegeLevel
}

type API struct {
	DB            *gorm.DB
	Config        *conf.Config
	WorkerPool    worker.Pool
	SDK           MateBotSDKGo.SDK
	EventChannels *map[*eventChannelKey]chan *eventNotification
}

func NewAPI(db *gorm.DB, config *conf.Config, client *MateBotSDKGo.SDK, wp worker.Pool) API {
	m := make(map[*eventChannelKey]chan *eventNotification)
	return API{
		DB:            db,
		Config:        config,
		WorkerPool:    wp,
		SDK:           *client,
		EventChannels: &m,
	}
}

type GenericResponse struct {
	Message string `json:"message"`
}

type simpleID struct {
	ID *uint64 `json:"id" echotools:"required"`
}

type namedObject struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type newMoneyRequest struct {
	Amount      *uint64 `json:"amount" echotools:"required"`
	Description *string `json:"description" echotools:"required;not empty"`
}

// Get the core user ID and the app's local user reference of the local authenticated user
// but without any validity checks (e.g., whether the user even exists at the core server)
func (a *API) getUnverifiedCoreID(c echo.Context) (uint64, *utilitymodels.LocalUser, error) {
	if context, err := middleware.GetSessionContext(c); err != nil {
		_ = c.JSON(500, GenericResponse{Message: "Unexpected failure"})
		return 0, nil, err
	} else {
		if !context.IsAuthenticated() {
			_ = c.JSON(401, GenericResponse{Message: "Unauthenticated"})
			return 0, nil, errors.New("unauthenticated")
		}
		u := context.GetUser()
		switch u.(type) {
		case *utilitymodels.LocalUser:
			var b models.CoreUser
			a.DB.Find(&b, "user_id = ?", u.(*utilitymodels.LocalUser).ID)
			if b.ID == 0 {
				_ = c.JSON(500, GenericResponse{Message: "Registered user for session not found"})
				return 0, nil, errors.New("session error")
			}
			return b.CoreID, u.(*utilitymodels.LocalUser), nil
		default:
			panic("invalid local user type")
		}
	}
}

// Get the active user's core user instance, if existent. Performs check on validity
// (e.g. whether the user is active, has a confirmed app alias, or minimal privilege level).
// If the function returns an error, the HTTP response has already been prepared.
func (a *API) getVerifiedCoreUser(c echo.Context, minimalLevel *MateBotSDKGo.PrivilegeLevel) (*MateBotSDKGo.User, *utilitymodels.LocalUser, error) {
	coreUserID, localUser, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return nil, nil, err
	}
	coreUser, err := a.SDK.GetVerifiedUser(coreUserID, minimalLevel)
	if err != nil {
		_ = c.JSON(400, GenericResponse{Message: err.Error()})
		return nil, nil, err
	}
	return coreUser, localUser, nil
}

// Return the local user ID for a given core user ID or nil if not found
func (a *API) findLocalUserID(coreUserID uint64) *uint64 {
	var user models.CoreUser
	if err := a.DB.Find(&user, "core_id = ?", coreUserID).Error; err != nil {
		return nil
	}
	u := user.UserID
	return &u
}

func (a *API) getAllCoreUserIDs() (*[]uint64, error) {
	var allUsers []models.CoreUser
	if err := a.DB.Find(&allUsers).Error; err != nil {
		return nil, err
	}
	coreUserMapping := make(map[uint64]bool)
	for _, u := range allUsers {
		coreUserMapping[u.CoreID] = true
	}
	_ = coreUserMapping
	coreUserIDs := make([]uint64, 0, len(allUsers))
	for u := range coreUserMapping {
		coreUserIDs = append(coreUserIDs, u)
	}
	return &coreUserIDs, nil
}

type vote struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Vote     bool   `json:"vote"`
}

func (a *API) convVote(v MateBotSDKGo.Vote) *vote {
	return &vote{
		UserID:   v.UserID,
		Username: v.Username,
		Vote:     v.Vote,
	}
}
