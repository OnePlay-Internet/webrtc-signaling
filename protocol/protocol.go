package protocol

import "github.com/OnePlay-Internet/signaling-module/webrtc/go/packet"

type Tenant interface {
	Send(*packet.UserResponse)
	Receive() *packet.UserRequest
	IsExited() bool
	Exit()
}

type OnTenantFunc func(token string, tent Tenant) error

type ProtocolHandler interface {
	OnTenant(fun OnTenantFunc)
}

type SignalingConfig struct {
	WebsocketPort int
	GrpcPort      int
}
