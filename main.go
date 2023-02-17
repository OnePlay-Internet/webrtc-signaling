package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pigeatgarlic/signaling/protocol"
	signalling "github.com/pigeatgarlic/signaling/signaling"
	"github.com/pigeatgarlic/signaling/validator"
	"github.com/pigeatgarlic/signaling/validator/oneplay"
	"github.com/pigeatgarlic/signaling/validator/thinkshare"
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
	valid := func() validator.Validator {
		switch schema {
		case "oneplay":
			return oneplay.NewOneplayValidator(validationUrl)
		case "thinkshare":
			return thinkshare.NewThinkshareValidator(validationUrl)
		default:
			return nil;
		}
	}()

	if valid == nil {
		fmt.Printf("unknown validator\n");
		return;
	}


	signalling.InitSignallingServer(&protocol.SignalingConfig{
		WebsocketPort: WebsocketPort,
		GrpcPort:      GrpcPort,
	}, valid)

	shutdown := make(chan bool)
	shutdown <- true
}
