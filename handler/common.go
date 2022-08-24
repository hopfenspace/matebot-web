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

type EventNotification struct {
	MinPrivilege MateBotSDKGo.PrivilegeLevel
	AllReceivers bool
	Receivers    *[]uint
	Data         any
}

type eventChannelKey struct {
	sessionID string
	userID    uint
	coreID    uint
	privilege MateBotSDKGo.PrivilegeLevel
}

type API struct {
	DB            *gorm.DB
	Config        *conf.Config
	WorkerPool    worker.Pool
	SDK           MateBotSDKGo.SDK
	EventChannels *map[*eventChannelKey]chan *EventNotification
}

func NewAPI(db *gorm.DB, config *conf.Config, client *MateBotSDKGo.SDK, wp worker.Pool) API {
	m := make(map[*eventChannelKey]chan *EventNotification)
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

func (a *API) getUser(c echo.Context) (uint, *utilitymodels.LocalUser, error) {
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
			return b.MateBotID, u.(*utilitymodels.LocalUser), nil
		default:
			panic("invalid local user type")
		}
	}
}

func (a *API) getUsers(c echo.Context) (*MateBotSDKGo.User, *utilitymodels.LocalUser, error) {
	coreUserID, user, err := a.getUser(c)
	if err != nil {
		return nil, nil, err
	}
	coreUser, err := a.SDK.GetUser(coreUserID, nil)
	if err != nil {
		return nil, nil, err
	}
	return coreUser, user, nil
}
