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
	volume = "volume"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "up"),
		"Was the last query of Gluster successful.",
		nil, nil,
	)

	volumesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "volumes_count"),
		"How many volumes were up at the last query.",
		nil, nil,
	)

	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "brick_count"),
		"Number of bricks at last query.",
		[]string{"volume"}, nil,
	)

	volumeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "volume_status"),
		"Status code of requested volume.",
		[]string{"volume"}, nil,
	)

	nodeSizeFreeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "node_size_free_bytes"),
		"Free bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeSizeTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, volume, "node_size_total_bytes"),
		"Total bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)
)

func ScrapeGlobalVolumeStatus(volumeStrings []string, allVolumes string, ch chan<- prometheus.Metric) error {
	// Collect metrics from volume info
	volumeInfo, err := ExecVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up = 0
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
		return err
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

	// Volume Count
	ch <- prometheus.MustNewConstMetric(
		volumesCount, prometheus.GaugeValue, float64(volumeInfo.VolInfo.Volumes.Count),
	)

	// Volume Status and Brick Count
	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if volumeStrings[0] == allVolumes || ContainsVolume(volumeStrings, volume.Name) {
			ch <- prometheus.MustNewConstMetric(
				brickCount, prometheus.GaugeValue, float64(volume.BrickCount), volume.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				volumeStatus, prometheus.GaugeValue, float64(volume.Status), volume.Name,
			)
		}
	}

	// Collect metrics from volume status all detail
	volumeStatusAll, err := ExecVolumeStatusAllDetail()
	if err != nil {
		return err
	}
	for _, vol := range volumeStatusAll.VolStatus.Volumes.Volume {
		for _, node := range vol.Node {
			if node.Status == 1 {
				ch <- prometheus.MustNewConstMetric(
					nodeSizeTotalBytes, prometheus.CounterValue, float64(node.SizeTotal), node.Hostname, node.Path, vol.VolName,
				)
				ch <- prometheus.MustNewConstMetric(
					nodeSizeFreeBytes, prometheus.CounterValue, float64(node.SizeFree), node.Hostname, node.Path, vol.VolName,
				)
			}
		}
	}

	return nil
}

// ExecVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeInfo() (VolumeInfoXML, error) {
	args := []string{"volume", "info"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolumeInfoXML{}, cmdErr
	}

	volumeInfo, err := VolumeInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// ExecVolumeStatusAllDetail executes "gluster volume status all detail" at the local machine
// returns VolumeStatusXML struct and error
func ExecVolumeStatusAllDetail() (VolumeStatusXML, error) {
	args := []string{"volume", "status", "all", "detail"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolumeStatusXML{}, cmdErr
	}
	volumeStatus, err := VolumeStatusAllDetailXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeStatus, err
	}
	return volumeStatus, nil
}

// VolumeInfoXMLUnmarshall unmarshalls bytes to VolumeInfoXML struct
func VolumeInfoXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeInfoXML, error) {
	var vol VolumeInfoXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}

// VolumeStatusAllDetailXMLUnmarshall reads bytes.buffer and returns unmarshalled xml
func VolumeStatusAllDetailXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeStatusXML, error) {
	var vol VolumeStatusXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}
