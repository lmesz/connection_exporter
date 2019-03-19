package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var (
	//tcp        0      0 127.0.0.1:44076         127.0.0.1:1445          TIME_WAIT   -
	pattern      = regexp.MustCompile(`(tcp|tcp6)\s+\d+\s+\d+\s+\d+.\d+.\d+.\d+.\d+:(\d+)\s+\d+.\d+.\d+.\d+.\d+:(\d+)\s+(\w+)`)
	defaultPort  = "9911"
	desc         = fmt.Sprintf("Port on which exporter should listen. Default: %s", defaultPort)
	portToListen = kingpin.Flag(
		"config.portToListen",
		desc,
	).Default(defaultPort).String()
	portsToWatch = kingpin.Flag(
		"config.portsToWatch",
		desc,
	).Default("3306,11211,7000,9160,8080,9911").String()
)

type connectionCollector struct {
	connectionMetric *prometheus.Desc
}

func newConnectionCollector() *connectionCollector {
	return &connectionCollector{
		connectionMetric: prometheus.NewDesc("connection_metric",
			"Shows different number of connections by port",
			[]string{"port", "state"}, nil,
		),
	}
}

func (collector *connectionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.connectionMetric
}

func (collector *connectionCollector) Collect(ch chan<- prometheus.Metric) {
	ports := strings.Split(*portsToWatch, ",")
	result := make(map[string]map[string]float64, len(ports))

	for port := range ports {
		result[ports[port]] = map[string]float64{
			"ESTABLISHED": 0.0,
			"SYN_SENT":    0.0,
			"SYN_RECV":    0.0,
			"FIN_WAIT1":   0.0,
			"FIN_WAIT2":   0.0,
			"TIME_WAIT":   0.0,
			"CLOSED":      0.0,
			"CLOSE_WAIT":  0.0,
			"LAST_ACK":    0.0,
			"LISTEN":      0.0,
			"CLOSING":     0.0,
			"UNKNOWN":     0.0,
		}
	}
	out, err := exec.Command("netstat", "-plantu").Output()
	if err != nil {
		log.Fatal(err)
	}

	splittedResult := strings.Split(string(out), "\n")

	for i := range splittedResult {
		line := splittedResult[i]
		if pattern.MatchString(line) {
			f := pattern.FindStringSubmatch(line)
			destination_port := f[3]
			state := f[4]
			if stringInSlice(destination_port, strings.Split(*portsToWatch, ",")) {
				result[destination_port][state]++
			}
		}
	}

	for source := range result {
		for state := range result[source] {
			ch <- prometheus.MustNewConstMetric(collector.connectionMetric, prometheus.CounterValue, result[source][state], source, state)
		}
	}
}

func main() {
	kingpin.Version(version.Print("syslog_ng_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	listenPort := fmt.Sprintf(":%s", *portToListen)

	prometheus.MustRegister(newConnectionCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port ", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
