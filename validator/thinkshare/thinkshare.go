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
	Queue []string `json:"queue"`
}
type TokenResp struct {
	Queue []string `json:"queue"`
	Route map[string]string `json:"route"`
}

func (val *ThinkshareValidator) Validate(queue []string) (map[string]string, []string) {
	buf,_ := json.Marshal(TokenReq{Queue: queue})
	resp,err := http.Post(val.url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return map[string]string{},queue
	}

	token_resp := &TokenResp{}
	data := make([]byte, 10000)
	n,err := resp.Body.Read(data)
	if err != nil {
		return map[string]string{},queue
	}

	json.Unmarshal(data[:n],token_resp);
	return token_resp.Route,token_resp.Queue
}
