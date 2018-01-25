package customer

import (
	"golang.org/x/net/context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/wgrpc"

	"github.com/hwgo/config"
	"github.com/hwgo/customer/proto"
)

type Client struct {
	*wgrpc.Client
	client proto.CustomerClient
}

func NewClient(tracer opentracing.Tracer, logger log.Factory) *Client {
	hostport := config.GetEndpoint("customer")

	ct := wgrpc.NewClient(hostport, tracer, logger)
	c := proto.NewCustomerClient(ct.Conn())

	return &Client{ct, c}
}

func (c *Client) Get(ctx context.Context) *proto.CustomerReply {
	defer c.Close()

	r, err := c.client.Get(ctx, &proto.CustomerRequest{Id: "760"})
	if err != nil {
		c.Logger().Info("could not greet: ", zap.Error(err))
		return nil
	} else {
		c.Logger().Info("Customer: ", zap.String("customer name", r.Name))
		return r
	}
}
