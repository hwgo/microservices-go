//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package customer

import (
	"golang.org/x/net/context"
	"time"

	// "github.com/hwgo/pher/delay"
	"github.com/hwgo/pher/wgrpc"

	// "github.com/hwgo/config"
	"github.com/hwgo/customer/proto"
)

type customerServer struct{}

func (s *customerServer) Get(context.Context, *proto.CustomerRequest) (*proto.CustomerReply, error) {
	// simulate RPC delay
	time.Sleep(7 * time.Millisecond)
	return &proto.CustomerReply{
			Id:       "218",
			Name:     "Tom",
			Location: "ChongQing",
		},
		nil
}

func NewServer(name string, hostPort string) *wgrpc.Server {
	s := wgrpc.NewServer(name, hostPort)
	proto.RegisterCustomerServer(s.GrpcServer, &customerServer{})
	return s
}
