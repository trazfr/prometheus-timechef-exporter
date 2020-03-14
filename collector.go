package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "timechef"
)

var (
	promDescSolde = prometheus.NewDesc(
		namespace+"_solde",
		"Remaining money on Timechef's account",
		[]string{"site"},
		nil)
)

type timechefCollector struct {
	fetcher *TimechefFetcher
}

func NewTimechefCollector(config *Config) prometheus.Collector {
	fetcher, err := NewTimecheFetcher(&http.Client{Timeout: config.Timeout},
		config.User,
		config.Password)
	if err != nil {
		log.Fatalf("Could not initialize Timechef: %s", err)
	}
	return &timechefCollector{
		fetcher: fetcher,
	}
}

func (t *timechefCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- promDescSolde
}

func (t *timechefCollector) Collect(ch chan<- prometheus.Metric) {
	timechefResults, err := t.fetcher.Fetch()
	if err == nil {
		ch <- prometheus.MustNewConstMetric(promDescSolde, prometheus.GaugeValue, timechefResults.Solde, timechefResults.Site)
	} else {
		log.Printf("Could not fetch the Timechef metrics; %s\n", err)
	}
}
