package thinkshare

import (
	"fmt"
	"testing"
)

func TestValidator(m *testing.T) {
	val := 	NewThinkshareValidator("http://localhost:4545/token/challenge/session");
	res,err := val.Validate("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWNpcGllbnQiOiI0MyIsImlzU2VydmVyIjoiRmFsc2UiLCJpZCI6IjIyIiwibmJmIjoxNjY0MTU0ODY5LCJleHAiOjE2NjQ0MTQwNjksImlhdCI6MTY2NDE1NDg2OX0.i8i73R7LGxfrFbyND2yI6xUHByA0eIEtMEA3iHT4jPo");
	if err != nil {
		m.Error(err);
	} else {
		fmt.Printf("%v\n",*res);
	}
}
