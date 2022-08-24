package handler

import (
	"encoding/json"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strings"
)

func (a *API) makeNotification(event MateBotSDKGo.Event) (*EventNotification, error) {
	// TODO
	return &EventNotification{
		MinPrivilege: MateBotSDKGo.External,
		AllReceivers: true,
		Receivers:    nil,
		Data:         event,
	}, nil
}

func (a *API) Callback(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return c.JSON(401, GenericResponse{Message: "Missing 'Authorization' header"})
	}
	authSlice := strings.Split(auth, " ")
	if len(authSlice) != 2 || authSlice[0] != "Bearer" {
		return c.JSON(400, GenericResponse{Message: "Badly formatted 'Authorization' header"})
	} else if a.Config.MateBot.CallbackSecret != nil && authSlice[1] != *a.Config.MateBot.CallbackSecret {
		return c.JSON(401, GenericResponse{Message: "Invalid secret"})
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: "Error while reading body"})
	}

	var events MateBotSDKGo.EventsNotification
	err = json.Unmarshal(body, &events)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: "Error while decoding json"})
	}

	notifications := make([]*EventNotification, len(events.Events))
	for i, event := range events.Events {
		notification, err := a.makeNotification(event)
		if err != nil {
			c.Logger().Error(err)
		} else {
			notifications[i] = notification
		}
	}

	for _, notificationChannel := range *a.EventChannels {
		notificationChannel := notificationChannel
		go func() {
			for _, notification := range notifications {
				notificationChannel <- notification
			}
		}()
	}

	return c.JSON(200, GenericResponse{Message: "OK"})
}
