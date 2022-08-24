package handler

import (
	"encoding/json"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strings"
)

func (a *API) makeNotification(event MateBotSDKGo.Event) (*eventWrapper, error) {
	// TODO
	return &eventWrapper{
		allUsers:     true,
		users:        nil,
		minPrivilege: MateBotSDKGo.External,
		notification: eventNotification{
			Type: event.Event,
			Data: event.Data,
		},
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

	notifications := make([]*eventWrapper, len(events.Events))
	for i, event := range events.Events {
		notification, err := a.makeNotification(event)
		if err != nil {
			c.Logger().Error(err)
		} else {
			notifications[i] = notification
		}
	}

	for identification, notificationChannel := range *a.EventChannels {
		identification := identification
		notificationChannel := notificationChannel
		go func() {
			for _, notification := range notifications {
				if notification.minPrivilege > identification.privilege {
					continue
				} else if !notification.allUsers && notification.users != nil {
					found := false
					for _, userID := range *notification.users {
						if userID == identification.coreID {
							found = true
						}
					}
					if !found {
						return
					}
				}
				notificationChannel <- &notification.notification
			}
		}()
	}

	return c.JSON(200, GenericResponse{Message: "OK"})
}
