package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pigeatgarlic/signaling/protocol"
	signalling "github.com/pigeatgarlic/signaling/signaling"
)

func main() {
	ValidationUrl := "https://auth.thinkmay.net/auth/validate";
	WebsocketPort := 8088;
	GrpcPort :=      8000;

	var err error
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--websocket" {
			WebsocketPort,err = strconv.Atoi(args[i+1])
		} else if arg == "--grpc" {
			GrpcPort,err = strconv.Atoi(args[i+1])
		} else if arg == "--validationurl" {
			ValidationUrl = args[i+1]
		} else if arg == "--help" {
			fmt.Printf("--engine |  encode engine ()\n")
			return
		}
	}

	if err != nil  {
		fmt.Printf("faile to parse argument: %s\n",err.Error());
		return;
	}

	signalling.InitSignallingServer(&protocol.SignalingConfig{
		WebsocketPort: WebsocketPort,
		GrpcPort:      GrpcPort,
		ValidationUrl: ValidationUrl,
	})

	shutdown := make(chan bool)
	shutdown <- true
}
