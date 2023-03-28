package signaling

import (
	"fmt"

	"github.com/thinkonmay/signaling-server/protocol"
	"github.com/thinkonmay/thinkremote-rtchub/signalling/gRPC/packet"
)

type Pair struct {
	A protocol.Tenant
	B protocol.Tenant
}

func (pair *Pair) handlePair() {
	pair.B.Send(&packet.SignalingMessage{
		Type: packet.SignalingType_tSTART,
		Sdp:  nil,
		Ice:  nil,
	})
	pair.A.Send(&packet.SignalingMessage{
		Type: packet.SignalingType_tSTART,
		Sdp:  nil,
		Ice:  nil,
	})


	stop := make(chan bool,2)
	go func() {
		for {
			msg := pair.B.Receive()
			if pair.A.IsExited() || pair.B.IsExited() {
				stop<-true
				return
			}
			if msg != nil {
				pair.A.Send(msg)
			}
		}
	}()
	go func() {
		for {
			msg := pair.A.Receive()
			if pair.A.IsExited() || pair.B.IsExited() {
				stop<-true
				return
			}
			if msg != nil {
				pair.B.Send(msg)
			}
		}
	}()
	go func() {
		<-stop
		fmt.Println("pair exited");
		if !pair.A.IsExited() {
			pair.A.Exit()
		}
		if !pair.B.IsExited() {
			pair.B.Exit()
		}
	}()

}
