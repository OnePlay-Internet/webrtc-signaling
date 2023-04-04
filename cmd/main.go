package main

import (
	// "github.com/pigeatgarlic/signaling/protocol"
	"github.com/thinkonmay/signaling-server"
	"github.com/thinkonmay/signaling-server/protocol"
	"github.com/thinkonmay/signaling-server/validator/sbvalidator"
	// "github.com/pigeatgarlic/signaling/validator"
)

const (
	default_signaling_grpc_port = 4000
	default_signaling_ws_port = 8080
	signaling_validate_url = "https://kczvtfaouddunjtxcemk.functions.supabase.co/signaling_authenticate"
	local_anon_key         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImtjenZ0ZmFvdWRkdW5qdHhjZW1rIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzk1NDc0MTcsImV4cCI6MTk5NTEyMzQxN30.dJqF_ipAx8NF_P__tsR-KkghVSc2McQo8B3MxeEup58"
)

func main() {



	signaling.InitSignallingServer(&protocol.SignalingConfig{
		WebsocketPort: default_signaling_ws_port,
		GrpcPort:      default_signaling_grpc_port,
	}, sbvalidator.NewSbValidator(signaling_validate_url,local_anon_key))



	shutdown := make(chan bool)
	shutdown <- true
}
