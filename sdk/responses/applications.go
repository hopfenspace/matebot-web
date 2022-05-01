package responses

type Application struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Created uint   `json:"created"`
}
