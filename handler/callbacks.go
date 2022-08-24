package handler

import (
	"encoding/json"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strings"
)

func (a *API) Callback(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return c.JSON(401, GenericResponse{Error: true, Message: "Missing 'Authorization' header"})
	}
	authSlice := strings.Split(auth, " ")
	if len(authSlice) != 2 || authSlice[0] != "Bearer" {
		return c.JSON(400, GenericResponse{Error: true, Message: "Badly formatted 'Authorization' header"})
	} else if a.Config.MateBot.CallbackSecret != nil && authSlice[1] != *a.Config.MateBot.CallbackSecret {
		return c.JSON(401, GenericResponse{Error: true, Message: "Invalid secret"})
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(400, GenericResponse{Error: true, Message: "Error while reading body"})
	}

	var events MateBotSDKGo.EventsNotification
	err = json.Unmarshal(body, &events)
	if err != nil {
		return c.JSON(400, GenericResponse{Error: true, Message: "Error while decoding json"})
	}

	return nil
}
