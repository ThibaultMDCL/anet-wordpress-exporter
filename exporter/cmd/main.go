package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ThibaultMDCL/anet-wordpress-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	UserAgent = "prometheus-wordpress-exporter"
)

var (
	portNum          = flag.Int("port", 11011, "The port to expose metrics to")
	configFile       = flag.String("config.file", "", "Configure which WordPress sites to monitor") // TODO : remettre une valeur par default a la fin des test pour forcé la de passé par fichier de conf
	monitorWordPress = flag.String("host", "", "Which host to monitor, with format <schema>://<host or FQDN>")
	useAuth          = flag.Bool("auth.basic", true, "Whether to use basic authentication (true|false)")
	authUsername     = flag.String("auth.user", "admin", "User to use with basic auth")
	authPassword     = flag.String("auth.pass", "admin", "Password to use with basic auth")
)

func init() {
	flag.Parse()

	if *configFile == "" {
		wp := utils.NewWordpress(
			"standalone",
			*monitorWordPress,
			UserAgent,
			*authUsername,
			*authPassword,
			*useAuth,
		)

		prometheus.MustRegister(utils.NewWordpressCollector(wp))
		return
	}

	config, err := utils.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	for _, target := range config.WordPresses {
		if !target.IsEnabled() {
			continue
		}

		wp := utils.NewWordpress(
			target.Name,
			target.URL,
			UserAgent,
			target.Username,
			target.ApplicationPassword,
			true,
		)

		prometheus.MustRegister(utils.NewWordpressCollector(wp))
		log.Printf("configured WordPress target: %s (%s)", target.Name, target.URL)
	}
}

func main() {

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Started WordPress exporter on port %d\n", *portNum)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}
