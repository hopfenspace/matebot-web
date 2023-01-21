package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/models"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utility"
	"strconv"
)

type simpleUser struct {
	UserID   *uint64 `json:"user_id"`
	CoreID   uint64  `json:"core_id"`
	Username string  `json:"username"`
}

type debtorUser struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
	Active   bool   `json:"active"`
}

type user struct {
	UserID     *uint64              `json:"user_id"`
	CoreID     uint64               `json:"core_id"`
	Username   string               `json:"username"`
	Balance    int64                `json:"balance"`
	Permission bool                 `json:"permission"`
	Active     bool                 `json:"active"`
	External   bool                 `json:"external"`
	VoucherID  interface{}          `json:"voucher_id"`
	Aliases    []MateBotSDKGo.Alias `json:"aliases"`
	Debtors    []debtorUser         `json:"debtors"`
	Created    uint64               `json:"created"`
	Modified   uint64               `json:"modified"`
}

type stateResponse struct {
	Message          string     `json:"message"`
	DetailsAvailable bool       `json:"details_available"`
	MinimalUser      simpleUser `json:"minimal_user"`
	User             *user      `json:"user"`
}

type aliasResponse struct {
	Message string              `json:"message"`
	Alias   *MateBotSDKGo.Alias `json:"alias"`
}

type listResponse struct {
	Message string       `json:"message"`
	Users   []simpleUser `json:"users"`
}

func (a *API) convUser(coreUser *MateBotSDKGo.User, localUserID *uint64, logger echo.Logger) *user {
	debtors := make([]debtorUser, 0)
	users, err := a.SDK.GetUsers(map[string]string{"active": "true", "voucher_id": strconv.FormatUint(coreUser.ID, 10), "community": "false"})
	if err != nil {
		logger.Error("Failed to lookup debtor users: ", err.Error())
	} else {
		for _, u := range users {
			debtors = append(debtors, debtorUser{
				UserID:   u.ID,
				Username: u.Name,
				Balance:  u.Balance,
				Active:   u.Active,
			})
		}
	}
	return &user{
		UserID:     localUserID,
		CoreID:     coreUser.ID,
		Username:   coreUser.Name,
		Balance:    coreUser.Balance,
		Permission: coreUser.Permission,
		Active:     coreUser.Active,
		External:   coreUser.External,
		VoucherID:  coreUser.VoucherID,
		Aliases:    coreUser.Aliases,
		Debtors:    debtors,
		Created:    coreUser.Created,
		Modified:   coreUser.Modified,
	}
}

func (a *API) Me(c echo.Context) error {
	coreUserID, localUser, err := a.getUnverifiedCoreID(c)
	if err != nil {
		return err
	}
	unverifiedCoreUser, err := a.SDK.GetUser(coreUserID, nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	localID := uint64(localUser.ID)
	if !a.SDK.IsUserConfirmed(unverifiedCoreUser) {
		return c.JSON(200, stateResponse{
			Message:          "OK",
			DetailsAvailable: false,
			User:             nil,
			MinimalUser: simpleUser{
				UserID:   &localID,
				CoreID:   unverifiedCoreUser.ID,
				Username: unverifiedCoreUser.Name,
			},
		})
	}
	coreUser := unverifiedCoreUser
	l := uint64(localUser.ID)
	return c.JSON(200, stateResponse{
		Message:          "OK",
		DetailsAvailable: true,
		User:             a.convUser(coreUser, &l, c.Logger()),
		MinimalUser: simpleUser{
			UserID:   &localID,
			CoreID:   coreUser.ID,
			Username: coreUser.Name,
		},
	})
}

type nameRequest struct {
	Name *string `json:"name" echotools:"required;not empty"`
}

func (a *API) ChangeUsername(c echo.Context) error {
	var r nameRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, localUser, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	user, err := a.SDK.SetUsername(coreUser.ID, *r.Name)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	localID := uint64(localUser.ID)
	return c.JSON(200, stateResponse{
		Message:          "OK",
		DetailsAvailable: true,
		User:             a.convUser(user, &localID, nil),
		MinimalUser: simpleUser{
			UserID:   &localID,
			CoreID:   user.ID,
			Username: user.Name,
		},
	})
}

type voucherRequest struct {
	Debtor *any `json:"debtor" echotools:"required"`
}

type voucherResponse struct {
	Message     string       `json:"message"`
	Debtor      debtorUser   `json:"debtor"`
	Voucher     *debtorUser  `json:"voucher"`
	Transaction *transaction `json:"transaction"`
}

// Only call this function with validated user IDs!
func (a *API) handleVouching(c echo.Context, voucher *uint64, issuer uint64) error {
	var r voucherRequest
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	var err error
	var voucherUpdate *MateBotSDKGo.VoucherUpdate
	switch (*r.Debtor).(type) {
	case string:
		if voucher == nil {
			voucherUpdate, err = a.SDK.SetVoucher((*r.Debtor).(string), nil, issuer)
		} else {
			voucherUpdate, err = a.SDK.SetVoucher((*r.Debtor).(string), *voucher, issuer)
		}
	case float64:
		if voucher == nil {
			voucherUpdate, err = a.SDK.SetVoucher(int((*r.Debtor).(float64)), nil, issuer)
		} else {
			voucherUpdate, err = a.SDK.SetVoucher(int((*r.Debtor).(float64)), *voucher, issuer)
		}
	default:
		return c.JSON(400, GenericResponse{Message: "Unknown JSON format for user"})
	}
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	d := debtorUser{
		UserID:   voucherUpdate.Debtor.ID,
		Username: voucherUpdate.Debtor.Name,
		Balance:  voucherUpdate.Debtor.Balance,
		Active:   voucherUpdate.Debtor.Active,
	}
	var v *debtorUser
	if voucherUpdate.Voucher != nil {
		v = &debtorUser{
			UserID:   voucherUpdate.Voucher.ID,
			Username: voucherUpdate.Voucher.Name,
			Balance:  voucherUpdate.Voucher.Balance,
			Active:   voucherUpdate.Voucher.Active,
		}
	}
	return c.JSON(200, voucherResponse{Message: "OK", Debtor: d, Voucher: v, Transaction: a.convTransaction(voucherUpdate.Transaction)})
}

func (a *API) StartVouching(c echo.Context) error {
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	return a.handleVouching(c, &coreUser.ID, coreUser.ID)
}

func (a *API) StopVouching(c echo.Context) error {
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return nil
	}
	return a.handleVouching(c, nil, coreUser.ID)
}

func (a *API) DropPrivileges(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."}) // TODO: Implement this
}

func (a *API) ConfirmAlias(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	alias, err := a.SDK.ConfirmAlias(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, aliasResponse{Message: "OK", Alias: alias})
}

type deletionResponse struct {
	Message string               `json:"message"`
	UserID  uint64               `json:"user_id"`
	Aliases []MateBotSDKGo.Alias `json:"aliases"`
}

func (a *API) DeleteAlias(c echo.Context) error {
	var r simpleID
	if err := utility.ValidateJsonForm(c, &r); err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	coreUser, _, err := a.getVerifiedCoreUser(c, nil)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	aliases, err := a.SDK.GetAliases(map[string]string{"id": strconv.Itoa(int(*r.ID))})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	if aliases[0].ApplicationID == a.SDK.GetThisApplicationID() {
		return c.JSON(400, GenericResponse{Message: "It's not possible to delete the currently used alias. Do you want to delete your account instead?"})
	}
	deletion, err := a.SDK.DeleteAlias(*r.ID, coreUser.ID)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	return c.JSON(200, deletionResponse{Message: "OK", UserID: deletion.UserID, Aliases: deletion.Aliases})
}

func (a *API) DeleteLocalAccount(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."}) // TODO: Implement this
}

func (a *API) DeleteFullAccount(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."}) // TODO: Implement this
}

func (a *API) ListUsers(c echo.Context) error {
	l := MateBotSDKGo.Vouched
	_, _, err := a.getVerifiedCoreUser(c, &l)
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	u, err := a.SDK.GetUsers(map[string]string{"active": "true", "community": "false", "alias_confirmed": "true"})
	if err != nil {
		return c.JSON(400, GenericResponse{Message: err.Error()})
	}
	users := make([]simpleUser, len(u))
	for i := range u {
		var b models.CoreUser
		a.DB.Find(&b, "core_id = ?", u[i].ID)
		if b.ID == 0 {
			users[i] = simpleUser{
				UserID:   nil,
				CoreID:   u[i].ID,
				Username: u[i].Name,
			}
		} else {
			users[i] = simpleUser{
				UserID:   &b.UserID,
				CoreID:   u[i].ID,
				Username: u[i].Name,
			}
		}
	}
	return c.JSON(200, listResponse{Message: "OK", Users: users})
}
