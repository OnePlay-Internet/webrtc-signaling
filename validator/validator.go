package validator

type Validator interface {
	Validate(queue []string) (map[string]string, []string)
}
