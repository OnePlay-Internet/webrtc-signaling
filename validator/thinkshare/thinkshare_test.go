package thinkshare

import (
	"fmt"
	"testing"
)

func TestValidator(m *testing.T) {
	val := 	NewThinkshareValidator("http://localhost:54321/functions/v1/","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0");
	res,err := val.Validate([]string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWNpcGllbnQiOiI0MyIsImlzU2VydmVyIjoiRmFsc2UiLCJpZCI6IjIyIiwibmJmIjoxNjY0MTU0ODY5LCJleHAiOjE2NjQ0MTQwNjksImlhdCI6MTY2NDE1NDg2OX0.i8i73R7LGxfrFbyND2yI6xUHByA0eIEtMEA3iHT4jPo"});

	if err != nil {
		m.Error(err);
	} else {
		fmt.Printf("%v\n",res);
	}
}
