package sbvalidator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thinkonmay/signaling-server/validator"
)

type SbValidator struct {
	url string
	key string
}

func NewSbValidator(url string, key string) validator.Validator {
	return &SbValidator{
		url: url,
		key: fmt.Sprintf("Bearer %s", key),
	}
}

type TokenResp struct {
	Queue []string         `json:"queue"`
	Pairs []validator.Pair `json:"pairs"`
}

func (val *SbValidator) Validate(queue []string) ([]validator.Pair, []string) {
	fmt.Printf("prevalidate: %d \n",len(queue))
	buf, _ := json.Marshal(queue)
	req, err := http.NewRequest("POST", val.url, bytes.NewBuffer(buf))
	if err != nil {
		fmt.Printf("%s\n",err.Error())
		return []validator.Pair{}, queue
	}

	req.Header.Set("Authorization", val.key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200{
		fmt.Printf("%s\n",err.Error())
		return []validator.Pair{}, queue
	}

	token_resp := &TokenResp{
		Pairs: []validator.Pair{},
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n",err.Error())
		return []validator.Pair{}, queue
	}

	err = json.Unmarshal(data, token_resp)
	if err != nil || resp.StatusCode != 200{
		fmt.Printf("%s\n",err.Error())
		return []validator.Pair{}, queue
	}
	fmt.Printf("aftervalidate: queue %d pairs: %d\n",len(token_resp.Queue),len(token_resp.Pairs))
	return token_resp.Pairs, token_resp.Queue
}
