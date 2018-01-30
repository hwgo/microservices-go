package echo

import (
	"net"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/tracing"
)

type Server struct {
	hostPort string
	addr     *net.UDPAddr
	tracer   opentracing.Tracer
	logger   log.Factory
}

func NewServer(hostPort string) *Server {
	logger := log.Service(ServiceName)
	tracer := tracing.Init(ServiceName, metrics.Namespace(ServiceName, nil), logger)

	//Resolving address
	udpAddr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		logger.Bg().Fatal("Error: ", zap.Error(err))
	}

	return &Server{
		hostPort: hostPort,
		addr:     udpAddr,
		tracer:   tracer,
		logger:   logger,
	}
}

// Run starts the Driver server
func (s *Server) Run() error {
	// Build listining connections
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	s.logger.Bg().Info("Create conn", zap.String("bind", s.addr.String()))

	// Interacting with one client at a time
	recvBuff := make([]byte, 1024)
	for {
		s.logger.Bg().Info("Ready to receive packets!")
		// Receiving a message
		rn, rmAddr, err := conn.ReadFromUDP(recvBuff)
		if err != nil {
			return err
		}

		s.logger.Bg().Info(
			"<<< Packet received",
			zap.String("from", rmAddr.String()),
			zap.String("data", string(recvBuff[:rn])),
		)

		// Sending the same message back to current client
		_, err = conn.WriteToUDP(recvBuff[:rn], rmAddr)
		if err != nil {
			return err
		}
		s.logger.Bg().Info(
			">>> Sent packet to",
			zap.String("to", rmAddr.String()),
		)
	}
}
