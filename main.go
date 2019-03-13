package main

import (
  "bufio"
  "net/http"
  //"os/exec"
  "os"
  "regexp"

  log "github.com/Sirupsen/logrus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/prometheus/client_golang/prometheus"
)

var (
  //tcp        0      0 127.0.0.1:44076         127.0.0.1:1445          TIME_WAIT   -
 pattern = regexp.MustCompile(`(tcp|tcp6)\s+\d+\s+\d+\s+\d+.\d+.\d+.\d+.\d+:(\d+)\s+\d+.\d+.\d+.\d+.\d+:(\d+)\s+(\w+)`)
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
  var result = make(map[string]map[string]float64)
  /*
  out, err := exec.Command("netstat -plantu").Output()
	if err != nil {
		log.Fatal(err)
	}
  */
  file, err := os.Open("test_data")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    if (pattern.MatchString(line)) {
      f := pattern.FindStringSubmatch(line)
      if _, ok := result[f[3]]; ok {
        if _, ok := result[f[3]][f[4]]; ok {
          result[f[3]][f[4]] = result[f[3]][f[4]] + 1
        } else {
          result[f[3]] = map[string]float64{}
          result[f[3]][f[4]] = 1
        }
      } else {
        result[f[3]] = map[string]float64{}
        result[f[3]][f[4]] = 1
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
  prometheus.MustRegister(newConnectionCollector())
  http.Handle("/metrics", promhttp.Handler())
  log.Info("Beginning to serve on port :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
