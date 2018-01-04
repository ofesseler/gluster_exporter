package collector

import (
	"bytes"
	"os/exec"
	"strings"
	"fmt"
	"time"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)


// TODO: Need Test
const (
	// Subsystem(s).
	mount = "mount"
)

type mountV struct {
	mountPoint 	string
	volume 		string
}

var (
	volumeWriteable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, mount, "mount_writable"),
		"Writes and deletes file in Volume and checks if it is writable",
		[]string{"volume", "mountpoint"}, nil)

	mountSuccessful = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, mount, "mount_successful"),
		"Checks if mountpoint exists, returns a bool value 0 or 1",
		[]string{"volume", "mountpoint"}, nil)
)

func ScrapeVolumeMountStatus(scrapeError *prometheus.CounterVec, ch chan<- prometheus.Metric) error {
	mountBuffer, execMountCheckErr := execMountCheck()
	if execMountCheckErr != nil {
		return execMountCheckErr
	} else {
		mounts, err := ParseMountOutput(mountBuffer.String())
		testMountResult := testMount(mounts, err)

		for _, mount := range mounts {
			ch <- prometheus.MustNewConstMetric(
				mountSuccessful, prometheus.GaugeValue, float64(testMountResult), mount.volume, mount.mountPoint,
			)
		}

		if err != nil {
			return err
		} else {
			for _, mount := range mounts {
				isWritable, err := execTouchOnVolumes(mount.mountPoint)
				if err != nil {
					scrapeError.WithLabelValues("collect.mount_status").Inc()
				}
				testWriteResult := testWritable(isWritable)
				ch <- prometheus.MustNewConstMetric(
					volumeWriteable, prometheus.GaugeValue, float64(testWriteResult), mount.volume, mount.mountPoint,
				)
			}
		}
	}

	return nil
}

func execMountCheck() (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	mountCmd := exec.Command("mount", "-t", "fuse.glusterfs")

	mountCmd.Stdout = stdoutBuffer
	err := mountCmd.Run()

	if err != nil {
		return stdoutBuffer, err
	}
	return stdoutBuffer, nil
}

// ParseMountOutput pares output of system execution 'mount'
func ParseMountOutput(mountBuffer string) ([]mountV, error) {
	mounts := make([]mountV, 0, 2)
	mountRows := strings.Split(mountBuffer, "\n")
	for _, row := range mountRows {
		trimmedRow := strings.TrimSpace(row)
		if len(row) > 3 {
			mountColumns := strings.Split(trimmedRow, " ")
			mounts = append(mounts, mountV{mountPoint: mountColumns[2], volume: mountColumns[0]})
		}
	}
	return mounts, nil
}

// Test is mount Successful or not
func testMount(mounts []mountV, err error) int {
	if mounts != nil && len(mounts) > 0 {
		return 1
	}
	return 0
}

// Test if mount Writable or not
func testWritable(isWritable bool) int {
	if isWritable {
		return 1
	}
	return 0
}

func execTouchOnVolumes(mountpoint string) (bool, error) {
	testFileName := fmt.Sprintf("%v/%v_%v", mountpoint, "gluster_mount.fixtures", time.Now())
	_, createErr := os.Create(testFileName)
	if createErr != nil {
		return false, createErr
	}
	removeErr := os.Remove(testFileName)
	if removeErr != nil {
		return false, removeErr
	}
	return true, nil
}
