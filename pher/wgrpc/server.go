package wgrpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/hwgo/pher/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/otgrpc"
)

type Server struct {
	Name     string
	Endpoint string

	LogFactory log.Factory
	Logger     log.Logger
	Tracer     opentracing.Tracer

	GrpcServer *grpc.Server
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.Endpoint)

	if err != nil {
		s.Logger.Fatal("Unable to start server", zap.Error(err))
		return err
	}

	s.Logger.Info("Starting", zap.String("address", "tcp://"+s.Endpoint))
	return s.GrpcServer.Serve(lis)
}

func NewServer(name string, hostPort string) *Server {
	logger := log.NewFactory(log.DefaultLogger.With(zap.String("service", name)))
	tracer := tracing.Init(name, metrics.Namespace(name, nil), logger)

	return &Server{
		Name:     name,
		Endpoint: hostPort,

		LogFactory: logger,
		Logger:     logger.Bg(),
		Tracer:     tracer,
		GrpcServer: newGrpcServer(tracer),
	}
}

func newGrpcServer(tracer opentracing.Tracer) *grpc.Server {
	th := otgrpc.NewTraceHandler(tracer, otgrpc.WithPayloadLogging())
	s := grpc.NewServer(grpc.StatsHandler(th))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return s
}
