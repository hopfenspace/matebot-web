package responses

type User struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Balance    int    `json:"balance"`
	Permission bool   `json:"permission"`
	Active     bool   `json:"active"`
	External   bool   `json:"external"`
	VoucherId  uint   `json:"voucher_id"`
	Aliases    []struct {
		Id            uint   `json:"id"`
		UserId        uint   `json:"user_id"`
		ApplicationId uint   `json:"application_id"`
		Username      string `json:"username"`
		Confirmed     bool   `json:"confirmed"`
	} `json:"aliases"`
	Created  uint `json:"created"`
	Modified uint `json:"modified"`
}
