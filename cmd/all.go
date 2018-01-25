package cmd

import "github.com/spf13/cobra"

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Starts all services",
	Long:  `Starts all services.`,
	Run: func(cmd *cobra.Command, args []string) {
		// logger.Info("Starting all services")
		go customerCmd.RunE(args)
		// go customerClientCmd.RunE(customerClientCmd, args)
		// go greeterCmd.RunE(greeterCmd, args)
		// go greeterClientCmd.RunE(greeterClientCmd, args)
		// go customerCmd.RunE(customerCmd, args)
		go driverCmd.RunE(driverCmd, args)
		// go routeCmd.RunE(routeCmd, args)
		frontendCmd.RunE(frontendCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(allCmd)
}
