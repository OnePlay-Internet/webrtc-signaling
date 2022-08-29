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

func (val *OneplayValidator) Validate(token string) (result *validator.TokenPair) {
	result = &validator.TokenPair{}
	resp,err := http.Get(fmt.Sprintf("%s/%s",val.url,token));
	if err != nil {
		fmt.Printf("validation failed, %s\n",err.Error());
		return nil
	}

	data := make([]byte,1000);
	n, err :=resp.Body.Read(data)
	if err != nil {
		fmt.Printf("validation failed, %s\n",err.Error());
		return nil
	}

	err = json.Unmarshal(data[:n],*result)
	if err != nil {
		fmt.Printf("validation failed, %s\n",err.Error());
		return nil
	}

	return result;
}