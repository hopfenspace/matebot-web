package requests

type CreateUser struct {
	Permission bool `json:"permission"`
	External   bool `json:"external"`
}

type CreateAlias struct {
	UserID        uint   `json:"user_id"`
	ApplicationID uint   `json:"application_id"`
	Username      string `json:"username"`
	Confirmed     bool   `json:"confirmed"`
}

type DeleteAlias struct {
	AliasID uint `json:"id"`
}

type UserFlagExternal struct {
	User     string `json:"user"`
	External bool   `json:"external"`
}

type UserFlagPermission struct {
	User       string `json:"user"`
	Permission bool   `json:"permission"`
}

type UserVoucher struct {
	Debtor  string `json:"debtor"`
	Voucher string `json:"voucher"`
}

type UserDisable struct {
	User string `json:"user"`
}
