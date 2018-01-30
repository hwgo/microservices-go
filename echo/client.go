package echo

import (
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
)

// Client is a remote client that implements driver.Interface
type Client struct {
	tracer opentracing.Tracer
	logger log.Factory
}

// NewClient creates a new driver.Client
func NewClient(tracer opentracing.Tracer, logger log.Factory) *Client {
	return &Client{
		tracer: tracer,
		logger: logger,
	}
}

func (c *Client) Hello(message string) {
	// c.logger.For(ctx).Info("Echo Hello", zap.String("message", message))

	// ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()
	// results, err := c.client.FindNearest(thrift.Wrap(ctx), location)
	// if err != nil {
	// 	return nil, err
	// }
	// return fromThrift(results), nil

	logger := c.logger.Bg()
	hostPort := viper.GetString("echo.server")

	// Resolving Address
	remoteAddr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		logger.Fatal("Cannot create echo client", zap.Error(err))
	}

	// Make a connection
	tmpAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}

	conn, err := net.DialUDP("udp", tmpAddr, remoteAddr)
	// Exit if some error occured
	if err != nil {
		logger.Fatal("Cannot dial to echo server", zap.Error(err))
	}
	defer conn.Close()

	// write a message to server
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		logger.Fatal("echo conn error", zap.Error(err))
	} else {
		logger.Info(">>> Packet sent to", zap.String("to", remoteAddr.String()))
	}

	// Receive response from server
	buf := make([]byte, 1024)
	rn, rmAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		logger.Fatal("echo conn read error", zap.Error(err))
	} else {
		logger.Info(
			"<<<  %d bytes received",
			zap.String("from", rmAddr.String()),
			zap.String("data", string(buf[:rn])),
		)
	}
}
