package main

import (
	"flag"
	"log"
	"strconv"
)

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

// loadConfig resolves the configuration in the following order, from highest to
// lowest priority:
//   - command-line flags
//   - environment variables
//   - defaults from NewConfig
//
// The flag set, arguments and environment lookup are passed in so that the
// resolution can be tested without touching the process globals.
func loadConfig(fs *flag.FlagSet, args []string, lookupEnv func(string) (string, bool)) (Config, error) {
	config := NewConfig()

	envString(lookupEnv, "REMOTE", &config.Remote)
	envString(lookupEnv, "SERVERQUERY_USER", &config.User)
	envString(lookupEnv, "SERVERQUERY_PASSWORD", &config.Password)
	envString(lookupEnv, "LISTEN_ADDRESS", &config.ListenAddr)
	envBool(lookupEnv, "ENABLE_CHANNEL_METRICS", &config.EnableChannelMetrics)
	envBool(lookupEnv, "IGNORE_FLOOD_LIMITS", &config.IgnoreFloodLimits)

	// The already resolved values act as the flag defaults, so an unset flag
	// keeps the environment value while an explicit flag overrides it.
	fs.StringVar(&config.Remote, "remote", config.Remote, "Remote address of server query port.")
	fs.StringVar(&config.User, "user", config.User, "The serverquery user of the ts3exporter.")
	fs.StringVar(&config.Password, "password", config.Password, "The password for the serverquery user.")
	fs.StringVar(&config.ListenAddr, "listen", config.ListenAddr, "Listen address of the exporter.")
	fs.BoolVar(&config.EnableChannelMetrics, "enablechannelmetrics", config.EnableChannelMetrics, "Enables the channel collector.")
	fs.BoolVar(&config.IgnoreFloodLimits, "ignorefloodlimits", config.IgnoreFloodLimits, "Disable the server query flood limiter. Use this only if your exporter is whitelisted in the query_ip_whitelist.txt file.")

	if err := fs.Parse(args); err != nil {
		return config, err
	}

	return config, nil
}

func envString(lookupEnv func(string) (string, bool), key string, target *string) {
	if value, found := lookupEnv(key); found {
		*target = value
	}
}

func envBool(lookupEnv func(string) (string, bool), key string, target *bool) {
	value, found := lookupEnv(key)
	if !found {
		return
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("ignoring invalid boolean value %q for %s\n", value, key)
		return
	}

	*target = parsed
}
