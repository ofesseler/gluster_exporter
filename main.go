// Copyright 2015 Oliver Fesseler
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Gluster exporter, exports metrics from gluster commandline tool.
package main

import (
	"flag"
	"net/http"

	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"os"
	"strings"
)

const (
	// GlusterCmd is the default path to gluster binary
	GlusterCmd = "/usr/sbin/gluster"
	namespace  = "gluster"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Gluster successful.",
		nil, nil,
	)

	volumesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volumes_count"),
		"How many volumes were up at the last query.",
		nil, nil,
	)

	volumeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_status"),
		"Status code of requested volume.",
		[]string{"volume"}, nil,
	)

	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_count"),
		"Number of bricks at last query.",
		[]string{"volume"}, nil,
	)

	peersConnected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peers_connected"),
		"Is peer connected to gluster cluster.",
		nil, nil,
	)
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	hostname string
	path     string
	volumes  []string
}

// Describe all the metrics exported by Gluster exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- volumeStatus
	ch <- volumesCount
	ch <- brickCount
	ch <- peersConnected
}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Collect metrics from volume info
	volumeInfo, err := ExecVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		log.Errorf("couldn't parse xml volume info: %v", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}

	// use OpErrno as indicator for up
	if volumeInfo.OpErrno != 0 {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		volumesCount, prometheus.GaugeValue, float64(volumeInfo.VolInfo.Volumes.Count),
	)

	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if e.volumes[0] == "_all" || ContainsVolume(e.volumes, volume.Name) {

			ch <- prometheus.MustNewConstMetric(
				brickCount, prometheus.GaugeValue, float64(volume.BrickCount), volume.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				volumeStatus, prometheus.GaugeValue, float64(volume.Status), volume.Name,
			)
		}
	}

	// reads gluster peer status
	peerStatus, err := ExecPeerStatus()
	if err != nil {
		log.Errorf("couldn't parse xml of peer status: %v", err)
	}
	count := 0
	for range peerStatus.Peer {
		count++
	}
	ch <- prometheus.MustNewConstMetric(
		peersConnected, prometheus.GaugeValue, float64(count),
	)

}

// ContainsVolume checks a slice if it cpntains a element
func ContainsVolume(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// NewExporter initialises exporter
func NewExporter(hostname, glusterExecPath, volumesString string) (*Exporter, error) {
	if len(glusterExecPath) < 1 {
		log.Fatalf("Gluster executable path is wrong: %v", glusterExecPath)
	}
	volumes := strings.Split(volumesString, ",")
	if len(volumes) < 1 {
		log.Warnf("No volumes given. Proceeding without volume information. Volumes: %v", volumesString)
	}

	return &Exporter{
		hostname: hostname,
		path:     glusterExecPath,
		volumes:  volumes,
	}, nil
}

func versionInfo() {
	fmt.Println(version.Print("gluster_exporter"))
	os.Exit(0)
}

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}

func main() {

	// commandline arguments
	var (
		glusterPath    = flag.String("gluster_executable_path", GlusterCmd, "Path to gluster executable.")
		metricPath     = flag.String("metrics-path", "/metrics", "URL Endpoint for metrics")
		listenAddress  = flag.String("listen-address", ":9189", "The address to listen on for HTTP requests.")
		showVersion    = flag.Bool("version", false, "Prints version information")
		glusterVolumes = flag.String("volumes", "_all", "Comma separated volume names: vol1,vol2,vol3. Default is '_all' to scrape all metrics")
	)
	flag.Parse()

	if *showVersion {
		versionInfo()
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	exporter, err := NewExporter(hostname, *glusterPath, *glusterVolumes)
	if err != nil {
		log.Errorf("Creating new Exporter went wrong, ... \n%v", err)
	}
	prometheus.MustRegister(exporter)

	log.Info("GlusterFS Metrics Exporter v", version.Version)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>GlusterFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>GlusterFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
