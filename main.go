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

// Glusterfs exorter currently scaping volume info
package main

import (
	"flag"
	"net/http"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"github.com/prometheus/common/log"
	"fmt"
	"github.com/prometheus/common/version"
	"os"
)

const (
	VERSION string = "0.1.0"
)

var (
	GLUSTER_CMD = "/usr/sbin/gluster"
)

type CliOutput struct {
	XMLName  xml.Name  `xml:"cliOutput"`
	OpRet    int         `xml:"opRet"`
	OpErrno  int       `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	VolInfo  VolInfo   `xml:"volInfo"`
}

type VolInfo struct {
	XMLName xml.Name  `xml:"volInfo"`
	Volumes Volumes   `xml:"volumes"`
}

type Volumes struct {
	XMLName xml.Name  `xml:"volumes"`
	Volume  []Volume   `xml:"volume"`
	Count   int         `xml:"count"`
}

type Volume struct {
	XMLName    xml.Name  `xml:"volume"`
	Name       string       `xml:"name"`
	Id         string         `xml:"id"`
	Status     int        `xml:"status"`
	StatusStr  string  `xml:"statusStr"`
	BrickCount int    `xml:"brickCount"`
	Bricks     []Brick    `xml:"bricks"`
	DistCount  int     `xml:"distCount"`
}

type Brick struct {
	Uuid      string       `xml:"brick>uuid"`
	Name      string       `xml:"brick>name"`
	HostUuid  string   `xml:"brick>hostUuid"`
	IsArbiter int     `xml:"brick>isArbiter"`
}

var (
	// Error number from GlusterFS
	errno = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:"glusterfs_errno",
			Help:"Error Number Glusterfs",
		},
		[]string{},
	)

	// creates a gauge of active nodes in glusterfs
	volume_count = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:"glusterfs_volume_count",
			Help:"Number of active glusterfs nodes",
		},
		[]string{},
	)

	// Count of bricks for gluster volume
	brick_count = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:"glusterfs_brick_count",
			Help:"Count of bricks for gluster volume",
		},
		[]string{"volume"},
	)

	// distribution count of bricks
	distribution_count = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:"glusterfs_nodes_active",
			Help:"distribution count of bricks",
		},
		[]string{"volume"},
	)
)

func init() {
	// register metric to prometheus's default registry
	prometheus.MustRegister(errno)
	prometheus.MustRegister(volume_count)
	prometheus.MustRegister(brick_count)
	prometheus.MustRegister(distribution_count)
}

func versionInfo() {
	fmt.Println("Gluster Exporter Version: ", VERSION)
	fmt.Println("Tested Gluster Version:   ", "3.8.5")
	fmt.Println("Go Version:               ", version.GoVersion)

	os.Exit(0)
}

func ExecGlusterCommand(arg ...string) *bytes.Buffer{
	stdoutBuffer := &bytes.Buffer{}
	glusterExec := exec.Command(GLUSTER_CMD, arg...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil  {
		log.Fatal(err)
	}
	return stdoutBuffer
}

// Unmarshall returned bytes to CliOutput struct
func infoUnmarshall(cmdOutBuff *bytes.Buffer) CliOutput {
	var vol CliOutput
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(b, &vol)
	return vol
}

func GlusterVolumeInfo() {
	// Execute gluster volume info
	stdOutbuff := ExecGlusterCommand("volume", "info")

	// Unmarshall returned bytes to CliOutput struct
	vol := infoUnmarshall(stdOutbuff)

	// set opErrno
	errno.WithLabelValues().Set(float64(vol.OpErrno))
	log.Debug("opErrno: %v", vol.OpErrno)

	// set volume count
	volume_count.WithLabelValues().Set(float64(vol.VolInfo.Volumes.Count))
	log.Debug("volume_count: %v", vol.VolInfo.Volumes.Count)

	// Volume based values
	for _, v := range vol.VolInfo.Volumes.Volume {
		// brick count with volume label
		brick_count.WithLabelValues(v.Name).Set(float64(v.BrickCount))
		log.Debug("opErrno: %v", vol.OpErrno)

		// distribution count with volume label
		distribution_count.WithLabelValues(v.Name).Set(float64(v.DistCount))
		log.Debug("opErrno: %v", vol.OpErrno)
	}
}

func glusterProfile(sec_int int) {
	// Gluster Profile


	// Get gluster volumes, then call gluster profile on every volume

	//  gluster volume profile gv_leoticket info cumulative --xml
	//cmd_profile := exec.Command("/usr/sbin/gluster", "volume", "profile", "gv_leoticket", "info", "cumulative", "--xml")
}

func main() {

	// commandline arguments
	var (
		metricPath = flag.String("metrics-path", "/metrics", "URL Endpoint for metrics")
		addr = flag.String("listen-address", ":9189", "The address to listen on for HTTP requests.")
		version_tag = flag.Bool("version", false, "Prints version information")
	)

	flag.Parse()

	if *version_tag {
		versionInfo()
	}

	log.Info("GlusterFS Metrics Exporter v", VERSION)

	// gluster volume info
	go GlusterVolumeInfo()

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
	log.Fatal(http.ListenAndServe(*addr, nil))
}
