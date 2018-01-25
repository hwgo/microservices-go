package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	metricsBackend string
)

var RootCmd = &cobra.Command{
	Use:   "hwgo",
	Short: "hwgo R.O.D. - A tracing demo application",
	Long:  `hwgo R.O.D. - A tracing demo application.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	RootCmd.PersistentFlags().StringVarP(&metricsBackend, "metrics", "m", "expvar", "Metrics backend (expvar|prometheus)")
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.BindPFlag("metrics", RootCmd.PersistentFlags().Lookup("metrics"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hwgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hwgo")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
