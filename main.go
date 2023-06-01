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
	config := NewConfig()
	setConfig(&config)

	flag.Parse()
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

func setConfig(config *Config) {
	if remote, found := os.LookupEnv("REMOTE"); found {
		config.Remote = remote
	} else {
		remoteFlag := flag.String("remote", "localhost:10011", "remote address of server query port")
		config.Remote = *remoteFlag
	}

	if user, found := os.LookupEnv("SERVERQUERY_USER"); found {
		config.User = user
	} else {
		userFlag := flag.String("user", "serveradmin", "the serverquery user of the ts3exporter")
		config.User = *userFlag
	}

	if pw, found := os.LookupEnv("SERVERQUERY_PASSWORD"); found {
		config.Password = pw
	} else {
		passwordFlag := flag.String("password", "", "The password for the serverquery user")
		config.Password = *passwordFlag
	}

	if listen, found := os.LookupEnv("LISTEN_ADDRESS"); found {
		config.ListenAddr = listen
	} else {
		listenFlag := flag.String("listen", "0.0.0.0:9189", "listen address of the exporter")
		config.ListenAddr = *listenFlag
	}

	if enableChannelMetrics, found := os.LookupEnv("ENABLE_CHANNEL_METRICS"); found {
		v, err := strconv.ParseBool(enableChannelMetrics)
		if err != nil {
			config.EnableChannelMetrics = false
		} else {
			config.EnableChannelMetrics = v
		}
	} else {
		enableChannelMetricsFlag := flag.Bool("enablechannelmetrics", false, "Enables the channel collector.")
		config.EnableChannelMetrics = *enableChannelMetricsFlag
	}

	if ignoreFloodLimits, found := os.LookupEnv("IGNORE_FLOOD_LIMITS"); found {
		v, err := strconv.ParseBool(ignoreFloodLimits)
		if err != nil {
			config.IgnoreFloodLimits = false
		} else {
			config.IgnoreFloodLimits = v
		}
	} else {
		ignoreFloodLimitsFlag := flag.Bool("ignorefloodlimits", false, "Disable the server query flood limiter. Use this only if your exporter is whitelisted in the query_ip_whitelist.txt file.")
		config.IgnoreFloodLimits = *ignoreFloodLimitsFlag
	}
}
