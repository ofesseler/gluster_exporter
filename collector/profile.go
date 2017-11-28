package collector

import (
	"bytes"
	"io/ioutil"
	"encoding/xml"

	"github.com/prometheus/common/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Subsystem(s).
	profile = "profile"
)

var (
	brickDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_duration"),
		"Time running volume brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataRead = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_data_read"),
		"Total amount of data read by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataWritten = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_data_written"),
		"Total amount of data written by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickFopHits = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_fop_hits"),
		"Total amount of file operation hits.",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyAvg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_fop_latency_avg"),
		"Average file operations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_fop_latency_min"),
		"Minimum file operations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, profile, "brick_fop_latency_max"),
		"Maximum file operations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)
)

func ScrapeProfileStatus(volumeStrings []string, allVolumes string, hostname string, scrapeError *prometheus.CounterVec, ch chan<- prometheus.Metric) error {
	// volumeInfo
	volumeInfo, err := ExecVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up = 0
	if err != nil {
		return err
	}

	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if volumeStrings[0] == allVolumes || ContainsVolume(volumeStrings, volume.Name) {
			volumeProfile, execVolProfileErr := ExecVolumeProfileGvInfoCumulative(volume.Name)
			if execVolProfileErr != nil {
				log.Errorf("Error while executing or marshalling gluster profile output: %v", execVolProfileErr)
				scrapeError.WithLabelValues("collect.profile_status").Inc()
			}

			for _, brick := range volumeProfile.Brick {

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
						brickFopLatencyAvg, prometheus.CounterValue, float64(fop.AvgLatency), volume.Name, brick.BrickName, fop.Name,
					)

					ch <- prometheus.MustNewConstMetric(
						brickFopLatencyMin, prometheus.CounterValue, float64(fop.MinLatency), volume.Name, brick.BrickName, fop.Name,
					)

					ch <- prometheus.MustNewConstMetric(
						brickFopLatencyMax, prometheus.CounterValue, float64(fop.MaxLatency), volume.Name, brick.BrickName, fop.Name,
					)
				}
			}
		}
	}

	return nil
}


// ExecVolumeProfileGvInfoCumulative executes "gluster volume profile {volume} info cumulative --xml" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeProfileGvInfoCumulative(volumeName string) (VolProfile, error) {
	args := []string{"volume", "profile"}
	args = append(args, volumeName)
	args = append(args, "info", "cumulative")
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolProfile{}, cmdErr
	}
	volumeProfile, err := VolumeProfileGvInfoCumulativeXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}

// VolumeProfileGvInfoCumulativeXMLUnmarshall unmarshalls cumulative profile of gluster volume profile
func VolumeProfileGvInfoCumulativeXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeProfileXML, error) {
	var vol VolumeProfileXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}
