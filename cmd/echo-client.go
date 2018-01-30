package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hwgo/echo"
	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/tracing"
)

var echoClientCmd *cobra.Command

func init() {
	name := "echo-client"

	echoClientCmd = &cobra.Command{
		Use:   name,
		Short: "Show.",
		Long:  `Show.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.Service(name)
			tracer := tracing.Init(name, metrics.Namespace(name, nil), logger)

			client := echo.NewClient(tracer, logger)
			client.Hello("echo")
			return nil
		},
	}

	RootCmd.AddCommand(echoClientCmd)
}
