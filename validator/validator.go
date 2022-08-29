package validator

type TokenPair struct {
	ServerToken string `json:"serverToken"`
	ClientToken string `json:"clientToken"`
}

type Validator interface {
	Validate(token string) *TokenPair
}