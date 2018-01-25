package wgrpc

import (
	"google.golang.org/grpc"

	"github.com/hwgo/pher/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/otgrpc"
)

type Client struct {
	tracer opentracing.Tracer
	logger log.Factory
	cc     *grpc.ClientConn
}

func NewClientWithTracing(name string, host string, port int) *Client {
	logger := log.NewFactory(log.DefaultLogger.With(zap.String("component", name)))
	tracer := tracing.Init(name, metrics.Namespace(name, nil), logger)

	return NewClient(host, port, tracer, logger)
}

func NewClient(host string, port int, tracer opentracing.Tracer, logger log.Factory) *Client {
	conn, err := newGrpcClientConn(HostPort(host, port), tracer)

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
