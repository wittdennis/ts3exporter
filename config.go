package main

type Config struct {
	Remote               string
	ListenAddr           string
	User                 string
	Password             string
	EnableChannelMetrics bool
	IgnoreFloodLimits    bool
}

func NewConfig() Config {
	config := Config{}
	config.Remote = "localhost:10011"
	config.ListenAddr = "0.0.0.0:9189"
	config.User = "serveradmin"
	config.Password = ""
	config.EnableChannelMetrics = false
	config.IgnoreFloodLimits = false

	return config
}
