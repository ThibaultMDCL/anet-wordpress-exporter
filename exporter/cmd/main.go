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

// func init() {
// }
func main() {
	flag.Parse()
	/*var config *utils.Config

	if *configFile != "" {
		var err error
		config, err = utils.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("unable to load configuration: %v", err)
		}
	}
	*/

	if *configFile != "" {
		if _, err := utils.LoadConfig(*configFile); err != nil {
			log.Fatalf("unable to load configuration: %v", err)
		}
	}

	wp := utils.NewWordpress(*monitorWordPress, UserAgent, *authUsername, *authPassword, *useAuth)
	prometheus.MustRegister(utils.NewWordpressCollector(wp))

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Started WordPress exporter for %s\n", *monitorWordPress)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}
