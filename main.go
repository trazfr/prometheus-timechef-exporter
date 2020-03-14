package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Args[0], "<config>")
		os.Exit(1)
	}
	config := NewConfig(os.Args[1])
	prometheus.MustRegister(NewTimechefCollector(config))

	http.Handle("/metrics", promhttp.Handler())
	log.Println(http.ListenAndServe(config.Listen, nil))
}
