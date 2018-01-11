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
	"net/http"

	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// GlusterCmd is the default path to gluster binary
	GlusterCmd = "/usr/sbin/gluster"
	namespace  = "gluster"
	allVolumes = "_all"
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

	nodeSizeFreeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_size_free_bytes"),
		"Free bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeSizeTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_size_bytes_total"),
		"Total bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeInodesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_inodes_total"),
		"Total inodes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeInodesFree = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_inodes_free"),
		"Free inodes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_available"),
		"Number of bricks available at last query.",
		[]string{"volume"}, nil,
	)

	brickDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_duration"),
		"Time running volume brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataRead = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_data_read"),
		"Total amount of data read by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataWritten = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_data_written"),
		"Total amount of data written by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickFopHits = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_hits_total"),
		"Total amount of file operation hits.",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyAvg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_avg"),
		"Average fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_min"),
		"Minimum fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_max"),
		"Maximum fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	peersConnected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peers_connected"),
		"Is peer connected to gluster cluster.",
		nil, nil,
	)

	healInfoFilesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "heal_info_files_count"),
		"File count of files out of sync, when calling 'gluster v heal VOLNAME info",
		[]string{"volume"}, nil)

	volumeWriteable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_writeable"),
		"Writes and deletes file in Volume and checks if it is writeable",
		[]string{"volume", "mountpoint"}, nil)

	mountSuccessful = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mount_successful"),
		"Checks if mountpoint exists, returns a bool value 0 or 1",
		[]string{"volume", "mountpoint"}, nil)

	quotaHardLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit"),
		"Quota hard limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaSoftLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_softlimit"),
		"Quota soft limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_used"),
		"Current data (bytes) used in a quota",
		[]string{"path", "volume"}, nil)

	quotaAvailable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_available"),
		"Current data (bytes) available in a quota",
		[]string{"path", "volume"}, nil)

	quotaSoftLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_softlimit_exceeded"),
		"Is the quota soft-limit exceeded",
		[]string{"path", "volume"}, nil)

	quotaHardLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit_exceeded"),
		"Is the quota hard-limit exceeded",
		[]string{"path", "volume"}, nil)
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	hostname string
	path     string
	volumes  []string
	profile  bool
	quota    bool
}

// Describe all the metrics exported by Gluster exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- volumeStatus
	ch <- volumesCount
	ch <- brickCount
	ch <- brickDuration
	ch <- brickDataRead
	ch <- brickDataWritten
	ch <- peersConnected
	ch <- nodeSizeFreeBytes
	ch <- nodeSizeTotalBytes
	ch <- brickFopHits
	ch <- brickFopLatencyAvg
	ch <- brickFopLatencyMin
	ch <- brickFopLatencyMax
	ch <- healInfoFilesCount
	ch <- volumeWriteable
	ch <- mountSuccessful
	ch <- quotaHardLimit
	ch <- quotaSoftLimit
	ch <- quotaUsed
	ch <- quotaAvailable
	ch <- quotaSoftLimitExceeded
	ch <- quotaHardLimitExceeded
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
		if e.volumes[0] == allVolumes || ContainsVolume(e.volumes, volume.Name) {

			ch <- prometheus.MustNewConstMetric(
				brickCount, prometheus.GaugeValue, float64(volume.BrickCount), volume.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				volumeStatus, prometheus.GaugeValue, float64(volume.Status), volume.Name,
			)
		}
	}

	// reads gluster peer status
	peerStatus, peerStatusErr := ExecPeerStatus()
	if peerStatusErr != nil {
		log.Errorf("couldn't parse xml of peer status: %v", peerStatusErr)
	}
	count := 0
	for range peerStatus.Peer {
		count++
	}
	ch <- prometheus.MustNewConstMetric(
		peersConnected, prometheus.GaugeValue, float64(count),
	)

	// reads profile info
	if e.profile {
		for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
			if e.volumes[0] == allVolumes || ContainsVolume(e.volumes, volume.Name) {
				volumeProfile, execVolProfileErr := ExecVolumeProfileGvInfoCumulative(volume.Name)
				if execVolProfileErr != nil {
					log.Errorf("Error while executing or marshalling gluster profile output: %v", execVolProfileErr)
				}
				for _, brick := range volumeProfile.Brick {
					if strings.HasPrefix(brick.BrickName, e.hostname) {
						ch <- prometheus.MustNewConstMetric(
							brickDuration, prometheus.CounterValue, float64(brick.CumulativeStats.Duration), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							brickDataRead, prometheus.CounterValue, float64(brick.CumulativeStats.TotalRead), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							brickDataWritten, prometheus.CounterValue, float64(brick.CumulativeStats.TotalWrite), volume.Name, brick.BrickName,
						)
						for _, fop := range brick.CumulativeStats.FopStats.Fop {
							ch <- prometheus.MustNewConstMetric(
								brickFopHits, prometheus.CounterValue, float64(fop.Hits), volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyAvg, prometheus.GaugeValue, fop.AvgLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyMin, prometheus.GaugeValue, fop.MinLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyMax, prometheus.GaugeValue, fop.MaxLatency, volume.Name, brick.BrickName, fop.Name,
							)
						}
					}
				}
			}
		}
	}

	// executes gluster status all detail
	volumeStatusAll, err := ExecVolumeStatusAllDetail()
	if err != nil {
		log.Errorf("couldn't parse xml of peer status: %v", err)
	}
	for _, vol := range volumeStatusAll.VolStatus.Volumes.Volume {
		for _, node := range vol.Node {
			ch <- prometheus.MustNewConstMetric(
				nodeSizeTotalBytes, prometheus.CounterValue, float64(node.SizeTotal), node.Hostname, node.Path, vol.VolName,
			)

			ch <- prometheus.MustNewConstMetric(
				nodeSizeFreeBytes, prometheus.GaugeValue, float64(node.SizeFree), node.Hostname, node.Path, vol.VolName,
			)
			ch <- prometheus.MustNewConstMetric(
				nodeInodesTotal, prometheus.CounterValue, float64(node.InodesTotal), node.Hostname, node.Path, vol.VolName,
			)

			ch <- prometheus.MustNewConstMetric(
				nodeInodesFree, prometheus.GaugeValue, float64(node.InodesFree), node.Hostname, node.Path, vol.VolName,
			)
		}
	}
	vols := e.volumes
	if vols[0] == allVolumes {
		log.Warn("no Volumes were given.")
		volumeList, volumeListErr := ExecVolumeList()
		if volumeListErr != nil {
			log.Error(volumeListErr)
		}
		vols = volumeList.Volume
	}

	for _, vol := range vols {
		filesCount, volumeHealErr := ExecVolumeHealInfo(vol)
		if volumeHealErr == nil {
			ch <- prometheus.MustNewConstMetric(
				healInfoFilesCount, prometheus.CounterValue, float64(filesCount), vol,
			)
		}
	}

	mountBuffer, execMountCheckErr := execMountCheck()
	if execMountCheckErr != nil {
		log.Error(execMountCheckErr)
	} else {
		mounts, err := parseMountOutput(mountBuffer.String())
		if err != nil {
			log.Error(err)
			if len(mounts) > 0 {
				for _, mount := range mounts {
					ch <- prometheus.MustNewConstMetric(
						mountSuccessful, prometheus.GaugeValue, float64(0), mount.volume, mount.mountPoint,
					)
				}
			}
		} else {
			for _, mount := range mounts {
				ch <- prometheus.MustNewConstMetric(
					mountSuccessful, prometheus.GaugeValue, float64(1), mount.volume, mount.mountPoint,
				)

				isWriteable, err := execTouchOnVolumes(mount.mountPoint)
				if err != nil {
					log.Error(err)
				}
				if isWriteable {
					ch <- prometheus.MustNewConstMetric(
						volumeWriteable, prometheus.GaugeValue, float64(1), mount.volume, mount.mountPoint,
					)
				} else {
					ch <- prometheus.MustNewConstMetric(
						volumeWriteable, prometheus.GaugeValue, float64(0), mount.volume, mount.mountPoint,
					)
				}
			}
		}
	}
	if e.quota {
		for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
			if e.volumes[0] == allVolumes || ContainsVolume(e.volumes, volume.Name) {
				volumeQuotaXML, err := ExecVolumeQuotaList(volume.Name)
				if err != nil {
					log.Error("Cannot create quota metrics if quotas are not enabled in your gluster server")
				} else {
					for _, limit := range volumeQuotaXML.VolQuota.QuotaLimits {
						ch <- prometheus.MustNewConstMetric(
							quotaHardLimit,
							prometheus.CounterValue,
							float64(limit.HardLimit),
							limit.Path,
							volume.Name,
						)

						ch <- prometheus.MustNewConstMetric(
							quotaSoftLimit,
							prometheus.CounterValue,
							float64(limit.SoftLimitValue),
							limit.Path,
							volume.Name,
						)
						ch <- prometheus.MustNewConstMetric(
							quotaUsed,
							prometheus.CounterValue,
							float64(limit.UsedSpace),
							limit.Path,
							volume.Name,
						)

						ch <- prometheus.MustNewConstMetric(
							quotaAvailable,
							prometheus.CounterValue,
							float64(limit.AvailSpace),
							limit.Path,
							volume.Name,
						)

						slExceeded := 0.0
						if limit.SlExceeded != "No" {
							slExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							quotaSoftLimitExceeded,
							prometheus.CounterValue,
							slExceeded,
							limit.Path,
							volume.Name,
						)

						hlExceeded := 0.0
						if limit.HlExceeded != "No" {
							hlExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							quotaHardLimitExceeded,
							prometheus.CounterValue,
							hlExceeded,
							limit.Path,
							volume.Name,
						)
					}
				}
			}
		}
	}
}

type mount struct {
	mountPoint string
	volume     string
}

// ParseMountOutput pares output of system execution 'mount'
func parseMountOutput(mountBuffer string) ([]mount, error) {
	mounts := make([]mount, 0, 2)
	mountRows := strings.Split(mountBuffer, "\n")
	for _, row := range mountRows {
		trimmedRow := strings.TrimSpace(row)
		if len(row) > 3 {
			mountColumns := strings.Split(trimmedRow, " ")
			mounts = append(mounts, mount{mountPoint: mountColumns[2], volume: mountColumns[0]})
		}
	}
	return mounts, nil
}

// ContainsVolume checks a slice if it contains an element
func ContainsVolume(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// NewExporter initialises exporter
func NewExporter(hostname, glusterExecPath, volumesString string, profile bool, quota bool) (*Exporter, error) {
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
		profile:  profile,
		quota:    quota,
	}, nil
}

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}

func main() {

	// commandline arguments
	var (
		metricsPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		listenAddress  = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9189").String()
		glusterPath    = kingpin.Flag("gluster.executable-path", "Path to gluster executable.").Default(GlusterCmd).String()
		glusterVolumes = kingpin.Flag("gluster.volumes", fmt.Sprintf("Comma separated volume names: vol1,vol2,vol3. Default is '%v' to scrape all metrics", allVolumes)).Default(allVolumes).String()
		profile        = kingpin.Flag("profile", "Enable gluster profiling reports.").Bool()
		quota          = kingpin.Flag("quota", "Enable gluster quota reports.").Bool()
		num            int
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("gluster_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting gluster_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("While trying to get Hostname error happened: %v", err)
	}
	exporter, err := NewExporter(hostname, *glusterPath, *glusterVolumes, *profile, *quota)
	if err != nil {
		log.Errorf("Creating new Exporter went wrong, ... \n%v", err)
	}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		num, err = w.Write([]byte(`<html>
			<head><title>GlusterFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>GlusterFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Fatal(num, err)
		}
	})

	log.Infoln("Listening on", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
