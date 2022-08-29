package main

import (
	"github.com/pigeatgarlic/signaling/protocol"
	signalling "github.com/pigeatgarlic/signaling/signaling"
)

func main() {
	shutdown := make(chan bool)
	signalling.InitSignallingServer(&protocol.SignalingConfig{
		WebsocketPort: 8088,
		GrpcPort:      8000,
		ValidationUrl: "https://auth.thinkmay.net/auth/validate",
	})
	shutdown <- true
}
