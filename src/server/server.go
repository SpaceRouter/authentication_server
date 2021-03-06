package server

import "github.com/spacerouter/authentication_server/config"

func Init() error {
	configs := config.GetConfig()
	r := NewRouter()
	return r.Run(configs.GetString("server.host") + ":" + configs.GetString("server.port"))
}
