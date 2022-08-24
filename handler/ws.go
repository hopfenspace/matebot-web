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
	_, _ = coreUser, user

	context, err := middleware.GetSessionContext(c)
	if err != nil {
		_ = c.JSON(500, GenericResponse{Error: true, Message: "Invalid session context"})
		return err
	}
	sessionID := context.GetSessionID()

	if _, exists := (*a.EventChannels)[*sessionID]; exists {
		c.Logger().Infof("WebSocket for session ID %s already exists", *sessionID)
		return c.JSON(400, GenericResponse{Error: true, Message: "WebSocket already set up"})
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				c.Logger().Error(err)
			}
		}(ws)
		for {
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return c.JSON(200, GenericResponse{Error: false, Message: "OK"})
}
