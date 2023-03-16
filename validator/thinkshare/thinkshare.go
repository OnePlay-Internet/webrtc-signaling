package thinkshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pigeatgarlic/signaling/validator"
)

type ThinkshareValidator struct {
	url string
	key string
}

func NewThinkshareValidator(url string,key string) validator.Validator {
	return &ThinkshareValidator{
		url: url,
		key: fmt.Sprintf("Bearer %s",key),
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
	req,err := http.NewRequest("POST",val.url,bytes.NewBuffer(buf))
	if err != nil {
		return map[string]string{},queue
	}

	req.Header.Set("Authorization",val.key)
	resp,err := http.DefaultClient.Do(req)
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
