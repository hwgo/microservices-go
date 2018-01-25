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

type customerServer struct {
	name   string
	server *wgrpc.Server
}

func (cs *customerServer) Get(ctx context.Context, cr *proto.CustomerRequest) (*proto.CustomerReply, error) {
	cs.server.LogFactory.For(ctx).Info("Get......Foo")
	// simulate RPC delay
	time.Sleep(7 * time.Millisecond)
	return &proto.CustomerReply{
			Id:       "218",
			Name:     "Tom",
			Location: "ChongQing",
		},
		nil
}

func NewServer(hostPort string) *wgrpc.Server {
	s := wgrpc.NewServer(ServiceName, hostPort)

	cs := &customerServer{
		name:   ServiceName,
		server: s,
	}

	proto.RegisterCustomerServer(s.GrpcServer, cs)
	return s
}
