package signaling

import (
	"github.com/thinkonmay/signaling-server/protocol"
	"github.com/thinkonmay/thinkremote-rtchub/signalling/gRPC/packet"
)

type Pair struct {
	client protocol.Tenant
	worker protocol.Tenant
}

func (pair *Pair) handlePair() {
	pair.worker.Send(&packet.SignalingMessage{
		Type: packet.SignalingType_START,
		Sdp:  nil,
		Ice:  nil,
	})
	go func() {
		for {
			if pair.client.IsExited() || pair.worker.IsExited() {
				return
			}
			pair.client.Send(pair.worker.Receive())
		}
	}()
	go func() {
		for {
			if pair.client.IsExited() || pair.worker.IsExited() {
				return
			}
			pair.worker.Send(pair.client.Receive())
		}
	}()
}
