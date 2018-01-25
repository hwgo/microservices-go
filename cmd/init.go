package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hwgo/config"
)

type ServerCommand struct {
	ServerName   string
	CobraCommand *cobra.Command
}

func (c *ServerCommand) AppendTo(root *cobra.Command) {
	var bind string
	var port int

	c.CobraCommand.Flags().StringVarP(&bind, "bind", "", "127.0.0.1",
		"interface to which the server will bind")
	c.CobraCommand.Flags().IntVarP(&port, "port", "p", 50052,
		"port on which the server will listen")

	c.SetOption("bind")
	c.SetOption("port")

	root.AddCommand(c.CobraCommand)
}

func (c *ServerCommand) SetOption(key string) {
	opt := c.ServerName + "." + key
	viper.BindEnv(opt)
	viper.BindPFlag(opt, c.CobraCommand.Flags().Lookup(key))
}

func (c *ServerCommand) RunE(args []string) error {
	return c.CobraCommand.RunE(c.CobraCommand, args)
}

func NewServerCommand(
	name string,
	run func(string) error,
) *ServerCommand {
	cmd := &cobra.Command{
		Use:   name,
		Short: "Show.",
		Long:  `Show.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint := config.GetEndpoint(name)
			return run(endpoint)
		},
	}

	return &ServerCommand{
		ServerName:   name,
		CobraCommand: cmd,
	}
}
