package grpc

import (
	"fmt"

	"github.com/thinkonmay/thinkremote-rtchub/signalling/gRPC/packet"
)

type GrpcTenant struct {
	exited bool
	client packet.Signaling_HandshakeServer

	pending chan *packet.SignalingMessage
}

func NewgRPCTenant(conn packet.Signaling_HandshakeServer) *GrpcTenant {
	ret := &GrpcTenant{
		pending: make(chan *packet.SignalingMessage, 5),
		client:  conn,
		exited:  false,
	}

	go func() {
		for {
			req, err := ret.client.Recv()
			if err != nil {
				ret.Exit()
			}
			ret.pending <- req
		}
	}()

	return ret
}

func (tenant *GrpcTenant) Send(pkt *packet.SignalingMessage) {
	if tenant.exited || pkt == nil {
		return
	}

	err := tenant.client.Send(pkt)
	if err != nil {
		tenant.Exit()
	}
}

func (tenant *GrpcTenant) Receive() *packet.SignalingMessage {
	if tenant.exited {
		return nil
	}

	return <-tenant.pending
}

func (tenant *GrpcTenant) IsExited() bool {
	return tenant.exited
}

func (tenant *GrpcTenant) Exit() {
	fmt.Printf("grpc tenant closed\n")
	tenant.pending <- nil
	tenant.exited = true
}

func (tenant *GrpcTenant) Peek() bool {
	return len(tenant.pending) > 0
}
