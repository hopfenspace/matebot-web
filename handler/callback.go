package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strconv"
	"strings"
)

func (a *API) makeNotification(event MateBotSDKGo.Event, logger echo.Logger) (*eventWrapper, error) {
	data := event.Data
	var receivers *[]uint64 = nil
	var minPrivilege *MateBotSDKGo.PrivilegeLevel
	minPrivilege = nil
	confirmedOnly := true
	var reply any

	switch event.Event {
	case MateBotSDKGo.ServerStarted:
		logger.Info("Core API server has been reported to be started")
		return nil, nil
	case MateBotSDKGo.AliasConfirmationRequested, MateBotSDKGo.AliasConfirmed:
		if data.App == nil || data.ID == nil || data.User == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		} else if *data.App == a.SDK.GetThisApplicationName() {
			return nil, nil
		}
		if aliases, err := a.SDK.GetAliases(map[string]string{"id": strconv.FormatUint(*data.ID, 10), "user_id": strconv.FormatUint(*data.User, 10)}); err != nil || len(aliases) != 1 {
			return nil, err
		} else {
			confirmedOnly = false
			receivers = &[]uint64{*data.User}
			reply = *aliases[0]
		}
	case MateBotSDKGo.CommunismCreated, MateBotSDKGo.CommunismUpdated, MateBotSDKGo.CommunismClosed:
		if data.ID == nil || data.Aborted == nil || data.Participants == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if communisms, err := a.SDK.GetCommunisms(map[string]string{"id": strconv.FormatUint(*data.ID, 10)}); err != nil || len(communisms) != 1 {
			return nil, err
		} else {
			m := MateBotSDKGo.Internal
			minPrivilege = &m
			r := make(map[uint64]bool)
			r[communisms[0].CreatorID] = true
			for _, p := range communisms[0].Participants {
				r[p.UserID] = true
			}
			rs := make([]uint64, 0)
			for i := range r {
				rs = append(rs, i)
			}
			receivers = &rs
			reply = a.convCommunism(communisms[0])
		}
	case MateBotSDKGo.PollCreated:
	case MateBotSDKGo.PollUpdated:
	case MateBotSDKGo.PollClosed:
	case MateBotSDKGo.RefundCreated:
	case MateBotSDKGo.RefundUpdated:
	case MateBotSDKGo.RefundClosed:
	case MateBotSDKGo.TransactionCreated:
	case MateBotSDKGo.VoucherUpdated:
	case MateBotSDKGo.UserSoftlyDeleted:
	case MateBotSDKGo.UserUpdated:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if user, err := a.SDK.GetUser(data.ID, nil); err != nil {
			return nil, err
		} else {
			receivers = &[]uint64{*data.ID}
			reply = user // TODO: convert the user to our schema
		}
	default:
		logger.Errorf("Unknown callback event type or type handler not implemented: '%s'", event.Event)
		return nil, nil
	}

	return &eventWrapper{
		receivers:     receivers,
		minPrivilege:  minPrivilege,
		confirmedOnly: confirmedOnly,
		notification: eventNotification{
			Type: event.Event,
			Data: reply,
		},
	}, nil
}

func sendNotification(notifications []*eventWrapper, identification *eventChannelKey, notificationChannel chan *eventNotification, _ echo.Logger) {
	for _, notification := range notifications {
		if notification == nil || notification.confirmedOnly && !identification.confirmed {
			// TODO: After being confirmed, the front-end should re-connect to the WebSocket to fix the `confirmed` attribute
			continue
		} else if notification.minPrivilege != nil && *notification.minPrivilege <= identification.privilege {
			notificationChannel <- &notification.notification
		} else if notification.receivers != nil {
			found := false
			for _, userID := range *notification.receivers {
				if userID == identification.coreID {
					found = true
				}
			}
			if !found {
				continue
			}
			notificationChannel <- &notification.notification
		}
	}
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
	if err = json.Unmarshal(body, &events); err != nil {
		return c.JSON(400, GenericResponse{Message: "Error while decoding json"})
	}

	notifications := make([]*eventWrapper, 0, len(events.Events))
	for _, event := range events.Events {
		if notification, err := a.makeNotification(event, c.Logger()); err != nil {
			c.Logger().Error(err)
		} else if notification != nil {
			notifications = append(notifications, notification)
		}
	}

	for identification, notificationChannel := range *a.EventChannels {
		go sendNotification(notifications, identification, notificationChannel, c.Logger())
	}

	return c.JSON(200, GenericResponse{Message: "OK"})
}
