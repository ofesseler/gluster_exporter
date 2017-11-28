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
	quota = "quota"
)

var (
	quotaHardLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_hardlimit"),
		"Quota hard limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaSoftLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_softlimit"),
		"Quota soft limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_used"),
		"Current data (bytes) used in a quota",
		[]string{"path", "volume"}, nil)

	quotaAvailable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_available"),
		"Current data (bytes) available in a quota",
		[]string{"path", "volume"}, nil)

	quotaSoftLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_softlimit_exceeded"),
		"Is the quota soft-limit exceeded",
		[]string{"path", "volume"}, nil)

	quotaHardLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, quota, "volume_quota_hardlimit_exceeded"),
		"Is the quota hard-limit exceeded",
		[]string{"path", "volume"}, nil)
)

func ScrapeQuotaStatus(volumeStrings []string, allVolumes string, scrapeError *prometheus.CounterVec, ch chan<- prometheus.Metric) error {
	// volumeInfo
	volumeInfo, err := ExecVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up = 0
	if err != nil {
		return err
	}

	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if volumeStrings[0] == allVolumes || ContainsVolume(volumeStrings, volume.Name) {

			if volumeQuota, err := ExecVolumeQuotaList(volume.Name); err != nil {
				log.Error("Cannot create quota metrics if quotas are not enabled in your Gluster Server")
				scrapeError.WithLabelValues("collect.quota_status").Inc()

			} else {
				for _, limit := range volumeQuota.VolQuota.QuotaLimits {
					ch <- prometheus.MustNewConstMetric(
						quotaHardLimit, prometheus.CounterValue, float64(limit.HardLimit), limit.Path, volume.Name,
					)

					ch <- prometheus.MustNewConstMetric(
						quotaSoftLimit, prometheus.CounterValue, float64(limit.SoftLimitValue), limit.Path, volume.Name,
					)

					ch <- prometheus.MustNewConstMetric(
						quotaUsed, prometheus.CounterValue, float64(limit.UsedSpace), limit.Path, volume.Name,
					)

					ch <- prometheus.MustNewConstMetric(
						quotaAvailable, prometheus.CounterValue, float64(limit.AvailSpace), limit.Path, volume.Name,
					)

					slExceeded := ExceededFunc(limit.SlExceeded)
					ch <- prometheus.MustNewConstMetric(
						quotaSoftLimitExceeded, prometheus.CounterValue, slExceeded, limit.Path, volume.Name,
					)

					hlExceeded := ExceededFunc(limit.HlExceeded)
					ch <- prometheus.MustNewConstMetric(
						quotaHardLimitExceeded, prometheus.CounterValue, hlExceeded, limit.Path, volume.Name,
					)
				}
			}
		}
	}

	return nil
}

// ExecVolumeQuotaList executes volume quota list on host system and processess input
// returns QuotaList structs and errors
func ExecVolumeQuotaList(volumeName string) (VolumeQuotaXML, error) {
	args := []string{"volume", "quota", volumeName, "list"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		// common error like "quota: No quota configured on volume {volume}"
		// return empty VolumeQuotaXML
		return VolumeQuotaXML{}, cmdErr
	}
	volumeQuota, err := VolumeQuotaListXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeQuota, err
	}
	return volumeQuota, nil
}

func VolumeQuotaListXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeQuotaXML, error) {
	var volQuotaXML VolumeQuotaXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return volQuotaXML, err
	}
	xml.Unmarshal(b, &volQuotaXML)
	return volQuotaXML, nil
}

func ExceededFunc(Exceeded string) float64 {
	if Exceeded != "No" {
		return 1.0
	}
	return 0.0
}
