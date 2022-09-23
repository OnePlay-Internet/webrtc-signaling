package thinkshare

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pigeatgarlic/signaling/validator"
)

type ThinkshareValidator struct {
	url string
}

func NewThinkshareValidator(url string) validator.Validator {
	return &ThinkshareValidator{
		url: url,
	}
}

type TokenReq struct {
	Token string `json:"token"`
}

func (val *ThinkshareValidator) Validate(token string) (result *validator.ValidationResult, err error) {
	result = &validator.ValidationResult{}

	buf, err := json.Marshal(TokenReq{Token: token})
	if err != nil {
		return
	}

	resp, err := http.Post(val.url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return
	}

	data := make([]byte, 1000)
	n, err := resp.Body.Read(data)
	if err != nil {
		return
	}

	err = json.Unmarshal(data[:n], result)
	if err != nil {
		return
	}
	return
}
