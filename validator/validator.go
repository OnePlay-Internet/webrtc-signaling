package validator

type ServerResult struct {
	ID int `json:"id"`
	ServerID int `json:"clientID"`
}

type ClientResult struct {
	ID int `json:"id"`
	ClientID int `json:"clientID"`
}

type Validator interface {
	ValidateServer(token string) (*ServerResult, error)
	ValidateClient(token string) (*ClientResult, error)
}