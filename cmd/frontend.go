package cmd

import (
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/tracing"

	"github.com/hwgo/frontend"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Starts Frontend service",
	Long:  `Starts Frontend service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.NewFactory(log.DefaultLogger.With(zap.String("service", "frontend")))

		server := frontend.NewServer(
			net.JoinHostPort(frontendOptions.serverInterface, strconv.Itoa(frontendOptions.serverPort)),
			tracing.Init("frontend", metrics.Namespace("frontend", nil), logger),
			logger,
		)
		return server.Run()
	},
}

var (
	frontendOptions struct {
		serverInterface string
		serverPort      int
	}
)

func init() {
	RootCmd.AddCommand(frontendCmd)

	frontendCmd.Flags().StringVarP(&frontendOptions.serverInterface, "bind", "", "127.0.0.1", "interface to which the frontend server will bind")
	frontendCmd.Flags().IntVarP(&frontendOptions.serverPort, "port", "p", 8080, "port on which the frontend server will listen")
}
