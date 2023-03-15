package validator

type ValidationResult struct {
	ID        int  `json:"id"`
	IsServer  bool `json:"isServer"`
	Recipient string `json:"recipient"`
}

type Validator interface {
	Validate(token string) (*ValidationResult, error)
}
