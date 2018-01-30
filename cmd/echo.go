package cmd

import (
	"github.com/spf13/viper"

	"github.com/hwgo/echo"
)

var echoCmd *ServerCommand

func init() {
	echoCmd = NewServerCommand(
		echo.ServiceName,
		func(endpoint string) error {
			hostPort := viper.GetString("echo.server")
			server := echo.NewServer(hostPort)
			return server.Run()
		},
	)
	echoCmd.AppendTo(RootCmd)
}
