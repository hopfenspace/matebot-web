package core

type Alias struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ApplicationID uint   `json:"application_id"`
	Username      string `json:"username"`
	Confirmed     bool   `json:"confirmed"`
}
