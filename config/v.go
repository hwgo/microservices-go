package config

import (
	"github.com/spf13/viper"

	"net"
	"strconv"
)

func GetEndpoint(prefix string) string {
	bind := viper.Get(prefix + ".bind")
	port := viper.Get(prefix + ".port")
	return net.JoinHostPort(bind.(string), strconv.Itoa(port.(int)))
}
