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
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
	Gs       *grpc.Server
}

// Run starts the Customer server
func (s *Server) Run() error {
	bg := s.logger.Bg()
	lis, err := net.Listen("tcp", s.hostPort)

	if err != nil {
		bg.Fatal("Unable to start server", zap.Error(err))
		return err
	}

	bg.Info("Starting", zap.String("address", "tcp://"+s.hostPort))
	return s.Gs.Serve(lis)
}

func NewServer(hostPort string, tracer opentracing.Tracer, logger log.Factory) *Server {
	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
		Gs:       newGrpcServer(tracer),
	}
}

func newGrpcServer(tracer opentracing.Tracer) *grpc.Server {
	th := otgrpc.NewTraceHandler(tracer, otgrpc.WithPayloadLogging())
	s := grpc.NewServer(grpc.StatsHandler(th))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return s
}

func NewServerWithTracing(name string, hostPort string) *Server {
	logger := log.NewFactory(log.DefaultLogger.With(zap.String("service", name)))
	tracer := tracing.Init(name, metrics.Namespace(name, nil), logger)

	return NewServer(hostPort, tracer, logger)
}
