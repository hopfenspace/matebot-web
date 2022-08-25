package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
	"golang.org/x/net/websocket"
)

func (a *API) WebSocket(c echo.Context) error {
	coreUser, user, err := a.getUsers(c)
	if err != nil {
		return err
	}

	context, err := middleware.GetSessionContext(c)
	if err != nil {
		_ = c.JSON(500, GenericResponse{Message: "Invalid session context"})
		return err
	}
	sessionID := context.GetSessionID()

	key := &eventChannelKey{
		sessionID: *sessionID,
		privilege: coreUser.Privilege(),
		coreID:    coreUser.ID,
	}

	if _, exists := (*a.EventChannels)[key]; exists {
		c.Logger().Infof("WebSocket for session ID %s already exists", *sessionID)
		return c.JSON(400, GenericResponse{Message: "WebSocket already set up"})
	}

	incoming := make(chan *eventNotification)
	(*a.EventChannels)[key] = incoming
	c.Logger().Infof("Registering new WebSocket for user ID %d ...", user.ID)

	websocket.Handler(func(ws *websocket.Conn) {
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				c.Logger().Error(err)
				delete(*a.EventChannels, key)
			}
		}(ws)
		for {
			data := <-incoming
			err = websocket.JSON.Send(ws, data)
			if err != nil {
				c.Logger().Error(err)
				delete(*a.EventChannels, key)
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
