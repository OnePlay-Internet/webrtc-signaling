package oneplay

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pigeatgarlic/signaling/validator"
)

type OneplayValidator struct {
	url string 
}

func NewOneplayValidator(url string) validator.Validator {
	return &OneplayValidator{
		url: url,
	}
}

func (val *OneplayValidator) ValidateServer(token string) (result *validator.ServerResult,err error) {
	result = &validator.ServerResult{}
	resp,err := http.Get(fmt.Sprintf("%s/%s",val.url,token));
	if err != nil {
		return 
	}

	data := make([]byte,1000);
	n, err :=resp.Body.Read(data)
	if err != nil {
		return 
	}

	err = json.Unmarshal(data[:n],result)
	if err != nil {
		fmt.Printf("validation failed, %s\n",err.Error());
		return 
	}

	return;
}
func (val *OneplayValidator) ValidateClient(token string) (result *validator.ClientResult, err error) {
	result = &validator.ClientResult{}
	resp,err := http.Get(fmt.Sprintf("%s/%s",val.url,token));
	if err != nil {
		return 
	}

	data := make([]byte,1000);
	n, err :=resp.Body.Read(data)
	if err != nil {
		return 
	}

	err = json.Unmarshal(data[:n],result)
	if err != nil {
		return 
	}
	return
}