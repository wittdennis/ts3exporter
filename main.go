package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/wittdennis/ts3exporter/pkg/collector"

	"github.com/wittdennis/ts3exporter/pkg/serverquery"
)

func main() {
	config := setConfig()

	c, err := serverquery.NewClient(config.Remote, config.User, config.Password, config.IgnoreFloodLimits)
	if err != nil {
		log.Fatalf("failed to init client %v\n", err)
	}

	internalMetrics := collector.NewExporterMetrics()
	seq := collector.SequentialCollector{collector.NewServerInfo(c, internalMetrics)}

	if config.EnableChannelMetrics {
		cInfo := collector.NewChannel(c, internalMetrics)
		seq = append(seq, cInfo)
	}

	prometheus.MustRegister(append(seq, collector.NewClient(c)))

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// setConfig loads configuration in the following order, from highest to lowest priority:
// - Environment variables
// - Command-line flags
// - Defaults from NewConfig()
func setConfig() Config {
	config := NewConfig()

	remote := config.Remote
	user := config.User
	password := config.Password
	listenAddr := config.ListenAddr
	enableChannelMetrics := config.EnableChannelMetrics
	ignoreFloodLimits := config.IgnoreFloodLimits

	flag.StringVar(&remote, "remote", remote, "Remote address of server query port.")
	flag.StringVar(&user, "user", user, "The serverquery user of the ts3exporter.")
	flag.StringVar(&password, "password", password, "The password for the serverquery user.")
	flag.StringVar(&listenAddr, "listen", listenAddr, "Listen address of the exporter.")
	flag.BoolVar(&enableChannelMetrics, "enablechannelmetrics", enableChannelMetrics, "Enables the channel collector.")
	flag.BoolVar(&ignoreFloodLimits, "ignorefloodlimits", ignoreFloodLimits, "Disable the server query flood limiter. Use this only if your exporter is whitelisted in the query_ip_whitelist.txt file.")

	flag.Parse()

	if env, found := os.LookupEnv("REMOTE"); found {
		config.Remote = env
	} else if isFlagPassed("remote") {
		config.Remote = remote
	}

	if env, found := os.LookupEnv("SERVERQUERY_USER"); found {
		config.User = env
	} else if isFlagPassed("user") {
		config.User = user
	}

	if env, found := os.LookupEnv("SERVERQUERY_PASSWORD"); found {
		config.Password = env
	} else if isFlagPassed("password") {
		config.Password = password
	}

	if env, found := os.LookupEnv("LISTEN_ADDRESS"); found {
		config.ListenAddr = env
	} else if isFlagPassed("listen") {
		config.ListenAddr = listenAddr
	}

	if env, found := os.LookupEnv("ENABLE_CHANNEL_METRICS"); found {
		if v, err := strconv.ParseBool(env); err == nil {
			config.EnableChannelMetrics = v
		}
	} else if isFlagPassed("enablechannelmetrics") {
		config.EnableChannelMetrics = enableChannelMetrics
	}

	if env, found := os.LookupEnv("IGNORE_FLOOD_LIMITS"); found {
		if v, err := strconv.ParseBool(env); err == nil {
			config.IgnoreFloodLimits = v
		}
	} else if isFlagPassed("ignorefloodlimits") {
		config.IgnoreFloodLimits = ignoreFloodLimits
	}

	return config
}
