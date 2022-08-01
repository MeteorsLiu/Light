package grpcclient

import (
	"crypto/tls"

	_grpc "github.com/MeteorsLiu/Light/gRPC-proto"
	"github.com/MeteorsLiu/Light/interfaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClient struct {
	client _grpc.LightClient
	config *interfaces.Config
}

func NewGRPCClient(config *interfaces.Config) (*GRPCClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	cfg := &tls.Config{}
	creds := credentials.NewTLS(cfg)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{}))
	conn, err := grpc.Dial(config.RemoteControllerAddr, opts...)
	client := _grpc.NewLightClient(conn)

	grpcClient := &GRPCClient{
		client: client,
	}

	grpcClient.Handshake()

	go grpcClient.Establish()

}

func (g *GRPCClient) Handshake() {

}

func (g *GRPCClient) Establish() {

}

func (g *GRPCClient) PacketHandler(upload *interfaces.UploadPayload) error {

}
