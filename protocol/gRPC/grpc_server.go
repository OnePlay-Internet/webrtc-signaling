package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/thinkonmay/signaling-server/protocol"
	"github.com/thinkonmay/thinkremote-rtchub/signalling/gRPC/packet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcServer struct {
	packet.UnimplementedSignalingServer
	grpcServer *grpc.Server
	fun        protocol.OnTenantFunc
}

func (server *GrpcServer) OnTenant(fun protocol.OnTenantFunc) {
	server.fun = fun
}

func InitSignallingServer(port int, options ...grpc.ServerOption) *GrpcServer {
	var ret GrpcServer
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}

	ret.grpcServer = grpc.NewServer(options...)
	packet.RegisterSignalingServer(ret.grpcServer, &ret)
	go ret.grpcServer.Serve(lis)
	return &ret
}

func (server *GrpcServer) Handshake(client packet.Signaling_HandshakeServer) error {
	md, ok := metadata.FromIncomingContext(client.Context())
	if !ok {
		return fmt.Errorf("Unauthorized")
	}

	token := md["authorization"]
	if token == nil {
		return fmt.Errorf("no authorize header")
	}

	tenant := NewgRPCTenant(client)
	err := server.fun(token[0], tenant)
	if err != nil {
		tenant.Exit()
	}

	for {
		if tenant.IsExited() {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
}
