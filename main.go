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
	"os/exec"

	"bytes"
	"fmt"
	"github.com/ofesseler/gluster_exporter/structs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"os"
	"strings"
)

const (
	namespace          = "gluster"
	VERSION     string = "0.1.0"
	GLUSTER_CMD        = "/usr/sbin/gluster"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Gluster successful.",
		[]string{"node"}, nil,
	)

	volumesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volumes_count"),
		"How many volumes were up at the last query.",
		[]string{"node"}, nil,
	)

	volumeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_status"),
		"Status code of requested volume.",
		[]string{"node", "volume"}, nil,
	)

	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_count"),
		"Number of bricks at last query.",
		[]string{"node", "volume"}, nil,
	)

	peerCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_count"),
		"Number of peers at last query.",
		[]string{"node"}, nil,
	)
)

type Exporter struct {
	hostname string
	path     string
	volumes  []string
}

// Describes all the metrics exported by Gluster exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- volumeStatus
	ch <- volumesCount
	ch <- brickCount
	ch <- peerCount
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Execute gluster volume info
	stdOutbuff := ExecGlusterCommand("volume", "info")
	// Unmarshall returned bytes to CliOutput struct
	vol, err := structs.VolumeInfoXmlUnmarshall(stdOutbuff)
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		log.Errorf("couldn't parse xml: %v", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0, e.hostname,
		)
	}

	// use OpErrno as indicator for up
	if vol.OpErrno != 0 {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0, e.hostname,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0, e.hostname,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		volumesCount, prometheus.GaugeValue, float64(vol.VolInfo.Volumes.Count), e.hostname,
	)

	for _, volume := range vol.VolInfo.Volumes.Volume {
		if volume.Name == "_all" || ContainsVolume(e.volumes, volume.Name) {

			ch <- prometheus.MustNewConstMetric(
				brickCount, prometheus.GaugeValue, float64(volume.BrickCount), e.hostname, volume.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				volumeStatus, prometheus.GaugeValue, float64(volume.Status), e.hostname, volume.Name,
			)
		}
	}

}

func ContainsVolume(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// comment
func NewExporter(hostname, glusterExecPath, volumes_string string) (*Exporter, error) {
	if len(glusterExecPath) < 1 {
		log.Fatalf("Gluster executable path is wrong: %v", glusterExecPath)
	}
	volumes := strings.Split(volumes_string, ",")
	if len(volumes) < 1 {
		log.Warnf("No volumes given. Proceeding without volume information. Volumes: %v", volumes_string)
	}

	return &Exporter{
		hostname: hostname,
		path:     glusterExecPath,
		volumes:  volumes,
	}, nil
}

func versionInfo() {
	fmt.Println("Gluster Exporter Version: ", VERSION)
	fmt.Println("Tested Gluster Version:   ", "3.8.5")
	fmt.Println("Go Version:               ", version.GoVersion)

	os.Exit(0)
}

func ExecGlusterCommand(arg ...string) *bytes.Buffer {
	stdoutBuffer := &bytes.Buffer{}
	arg_xml := append(arg, "--xml")
	glusterExec := exec.Command(GLUSTER_CMD, arg_xml...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		log.Fatal(err)
	}
	return stdoutBuffer
}

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}

func main() {

	// commandline arguments
	var (
		glusterPath    = flag.String("gluster_executable_path", GLUSTER_CMD, "Path to gluster executable.")
		metricPath     = flag.String("metrics-path", "/metrics", "URL Endpoint for metrics")
		listenAddress  = flag.String("listen-address", ":9189", "The address to listen on for HTTP requests.")
		showVersion    = flag.Bool("version", false, "Prints version information")
		glusterVolumes = flag.String("volumes", "_all", "Comma seperated volume names: vol1,vol2,vol3. Default is '_all' to scrape all metrics")
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

	log.Info("GlusterFS Metrics Exporter v", VERSION)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>GlusterFS Exporter v` + VERSION + `</title></head>
			<body>
			<h1>GlusterFS Exporter v` + VERSION + `</h1>
			<p><a href='` + *metricPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
