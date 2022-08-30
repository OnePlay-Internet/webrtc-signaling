package validator

type ValidationResult struct {
	ID        int  `json:"id"`
	IsServer  bool `json:"isServer"`
	Recipient bool `json:"recipient"`
}

type Validator interface {
	Validate(token string) (*ValidationResult, error)
}
