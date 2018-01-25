package wgrpc

import (
	"google.golang.org/grpc"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/otgrpc"
)

type Client struct {
	tracer opentracing.Tracer
	logger log.Factory
	cc     *grpc.ClientConn
}

func NewClient(hostport string, tracer opentracing.Tracer, logger log.Factory) *Client {
	conn, err := newGrpcClientConn(hostport, tracer)

	if err != nil {
		logger.Bg().Fatal("did not connect: ", zap.Error(err))
	}

	return &Client{
		tracer: tracer,
		logger: logger,
		cc:     conn,
	}
}

func newGrpcClientConn(hostport string, tracer opentracing.Tracer) (*grpc.ClientConn, error) {
	th := otgrpc.NewTraceHandler(tracer, otgrpc.WithPayloadLogging())
	return grpc.Dial(
		hostport,
		grpc.WithStatsHandler(th),
		grpc.WithInsecure(),
	)
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.cc
}

func (c *Client) Logger() log.Logger {
	return c.logger.Bg()
}

func (c *Client) Close() {
	c.cc.Close()
}

func (c *Client) LoggerFactory() log.Factory {
	return c.logger
}
