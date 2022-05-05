package core

type User struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Balance    int     `json:"balance"`
	Permission bool    `json:"permission"`
	Active     bool    `json:"active"`
	External   bool    `json:"external"`
	VoucherID  uint    `json:"voucher_id"`
	Aliases    []Alias `json:"aliases"`
	Created    int     `json:"created"`
	Modified   int     `json:"modified"`
}
