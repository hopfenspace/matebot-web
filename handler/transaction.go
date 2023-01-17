package handler

import (
	"fmt"
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
)

type transactionResponse struct {
	Message     string       `json:"message"`
	Transaction *transaction `json:"transaction"`
}

type transaction struct {
	Id        uint       `json:"id"`
	Sender    simpleUser `json:"sender"`
	Receiver  simpleUser `json:"receiver"`
	Amount    uint       `json:"amount"`
	Reason    *string    `json:"reason"`
	Timestamp uint       `json:"timestamp"`
}

type multiTransaction struct {
	BaseAmount   uint          `json:"base_amount"`
	TotalAmount  uint          `json:"total_amount"`
	Transactions []transaction `json:"transactions"`
	Timestamp    uint          `json:"timestamp"`
}

func (a *API) convTransaction(t *MateBotSDKGo.Transaction) *transaction {
	if t == nil {
		return nil
	}
	senderUserID := a.findLocalUserID(t.Sender.ID)
	receiverUserID := a.findLocalUserID(t.Receiver.ID)
	return &transaction{
		Id: t.ID,
		Sender: simpleUser{
			UserID:   senderUserID,
			CoreID:   t.Sender.ID,
			Username: t.Sender.Name,
		},
		Receiver: simpleUser{
			UserID:   receiverUserID,
			CoreID:   t.Receiver.ID,
			Username: t.Receiver.Name,
		},
		Amount:    t.Amount,
		Reason:    t.Reason,
		Timestamp: t.Timestamp,
	}
}

type sendTransactionRequest struct {
	Receiver any     `json:"receiver"`
	Amount   *uint   `json:"amount" echotools:"required"`
	Reason   *string `json:"reason" echotools:"required;not empty"`
}

func (a *API) SendTransaction(c echo.Context) error {
	var r sendTransactionRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return err
	}

	switch r.Receiver.(type) {
	case float64:
		transaction, err := a.SDK.SendTransaction(coreID, int(r.Receiver.(float64)), *r.Amount, *r.Reason)
		if err != nil {
			return c.JSON(400, GenericResponse{Message: err.Error()})
		}
		return c.JSON(200, transactionResponse{Message: "OK", Transaction: a.convTransaction(transaction)})
	case string:
		transaction, err := a.SDK.SendTransaction(coreID, r.Receiver.(string), *r.Amount, *r.Reason)
		if err != nil {
			return c.JSON(400, GenericResponse{Message: err.Error()})
		}
		return c.JSON(200, transactionResponse{Message: "OK", Transaction: a.convTransaction(transaction)})
	default:
		return c.JSON(400, GenericResponse{Message: fmt.Sprintf("Invalid data type %T", r.Receiver)})
	}
}

type consumeTransactionRequest struct {
	Amount     *uint   `json:"amount" echotools:"required"`
	Consumable *string `json:"consumable" echotools:"required;not empty"`
}

func (a *API) ConsumeTransaction(c echo.Context) error {
	var r consumeTransactionRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	_, _, _ := a.getCoreUser(c)
	coreID, _, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return err
	}
	t, err := a.SDK.ConsumeTransaction(coreID, *r.Amount, *r.Consumable)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, transactionResponse{Message: "OK", Transaction: a.convTransaction(t)})
}
