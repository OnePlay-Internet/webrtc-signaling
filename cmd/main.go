package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pigeatgarlic/signaling/protocol"
	"github.com/pigeatgarlic/signaling"
	"github.com/pigeatgarlic/signaling/validator"
)

func main() {
	validationUrl := os.Getenv("VALIDATION_URL")
	schema := os.Getenv("SCHEMA")

	WebsocketPort := 8088
	GrpcPort := 8000

	var err error
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--websocket" {
			WebsocketPort, err = strconv.Atoi(args[i+1])
		} else if arg == "--grpc" {
			GrpcPort, err = strconv.Atoi(args[i+1])
		} else if arg == "--validationurl" {
			validationUrl = args[i+1]
		} else if arg == "--schema" {
			schema = args[i+1]
		} else if arg == "--help" {
			fmt.Printf("--engine |  encode engine ()\n")
			return
		}
	}

	if err != nil {
		fmt.Printf("faile to parse argument: %s\n", err.Error())
		return
	}

	if schema == "" {
		schema = "thinkshare"
	}





	signaling.InitSignallingServer(&protocol.SignalingConfig{
		WebsocketPort: WebsocketPort,
		GrpcPort:      GrpcPort,
	}, nil)

	shutdown := make(chan bool)
	shutdown <- true
}
