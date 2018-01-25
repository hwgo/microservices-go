package cmd

import (
	"github.com/hwgo/customer"
)

var customerCmd *ServerCommand

func init() {
	name := "customer"
	customerCmd = NewServerCommand(
		name,
		func(endpoint string) error {
			server := customer.NewServer(name, endpoint)
			return server.Run()
		},
	)
	customerCmd.AppendTo(RootCmd)
}
