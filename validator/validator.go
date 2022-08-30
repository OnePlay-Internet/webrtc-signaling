package validator

type ValidationResult struct {
	ID        int  `json:"id"`
	IsServer  bool `json:"isServer"`
	Recipient int `json:"recipient"`
}

type Validator interface {
	Validate(token string) (*ValidationResult, error)
}
