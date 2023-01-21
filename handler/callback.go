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

type smallVoucherUpdate struct {
	Debtor      debtorUser   `json:"debtor"`
	Voucher     *uint64      `json:"voucher"`
	Transaction *transaction `json:"transaction"`
}

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
	case MateBotSDKGo.PollCreated, MateBotSDKGo.PollUpdated, MateBotSDKGo.PollClosed:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if polls, err := a.SDK.GetPolls(map[string]string{"id": strconv.FormatUint(*data.ID, 10)}); err != nil || len(polls) != 1 {
			return nil, err
		} else {
			m := MateBotSDKGo.Permitted
			minPrivilege = &m
			receivers = &[]uint64{polls[0].CreatorID, polls[0].User.ID}
			reply = a.convPoll(polls[0])
		}
	case MateBotSDKGo.RefundCreated, MateBotSDKGo.RefundUpdated, MateBotSDKGo.RefundClosed:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if refunds, err := a.SDK.GetRefunds(map[string]string{"id": strconv.FormatUint(*data.ID, 10)}); err != nil || len(refunds) != 1 {
			return nil, err
		} else {
			m := MateBotSDKGo.Permitted
			minPrivilege = &m
			receivers = &[]uint64{refunds[0].Creator.ID}
			reply = a.convRefund(refunds[0])
		}
	case MateBotSDKGo.TransactionCreated:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if transactions, err := a.SDK.GetTransactions(map[string]string{"id": strconv.FormatUint(*data.ID, 10)}); err != nil || len(transactions) != 1 {
			return nil, err
		} else {
			m := MateBotSDKGo.Permitted
			minPrivilege = &m
			receivers = &[]uint64{transactions[0].Sender.ID, transactions[0].Receiver.ID}
			reply = a.convTransaction(transactions[0])
		}
	case MateBotSDKGo.VoucherUpdated:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if user, err := a.SDK.GetUser(*data.ID, nil); err != nil {
			return nil, err
		} else {
			d := debtorUser{
				UserID:   user.ID,
				Username: user.Name,
				Balance:  user.Balance,
				Active:   user.Active,
			}
			var t *transaction = nil
			if data.Transaction != nil {
				if ts, err := a.SDK.GetTransactions(map[string]string{"id": strconv.FormatUint(*data.Transaction, 10)}); err == nil && len(ts) != 1 {
					t = a.convTransaction(ts[0])
				}
			}
			v := data.Voucher
			if data.Voucher != nil {
				if *data.Voucher != *user.VoucherID {
					logger.Error("Voucher as determined by new lookup differs from voucher in callback event data (using newly looked up data)!")
					v = user.VoucherID
				}
				receivers = &[]uint64{*data.ID, *data.Voucher}
			} else {
				receivers = &[]uint64{*data.ID}
			}
			reply = smallVoucherUpdate{Debtor: d, Voucher: v, Transaction: t}
		}
	case MateBotSDKGo.UserSoftlyDeleted:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if user, err := a.SDK.GetUser(data.ID, nil); err != nil {
			return nil, err
		} else {
			receivers = &[]uint64{*data.ID}
			reply = a.convUser(user, a.findLocalUserID(user.ID), logger)
		}
	case MateBotSDKGo.UserUpdated:
		if data.ID == nil {
			return nil, errors.New(fmt.Sprintf("Invalid incoming callback event %s, because expected fields were unset", event.Event))
		}
		if user, err := a.SDK.GetUser(data.ID, nil); err != nil {
			return nil, err
		} else if user.Active {
			return nil, errors.New(fmt.Sprintf("User %d seems to be active, even though a callback event USER_SOFTLY_DELETED has been received", *data.ID))
		} else {
			receivers = &[]uint64{*data.ID}
			reply = a.convUser(user, a.findLocalUserID(user.ID), logger)
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
