package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/wittdennis/ts3exporter/pkg/collector"

	"github.com/wittdennis/ts3exporter/pkg/serverquery"
)

func main() {
	config, err := loadConfig(flag.CommandLine, os.Args[1:], os.LookupEnv)
	if err != nil {
		log.Fatalf("failed to load config %v\n", err)
	}

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
