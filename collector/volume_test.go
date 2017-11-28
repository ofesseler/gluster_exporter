package collector

import (
	"testing"
	"io/ioutil"
	"bytes"
	"fmt"
)

func TestVolumeInfoXMLUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/gluster_volume_info.xml")
	if err != nil {
		t.Fatal(err)
	}

	// Convert into bytes.buffer
	contentBuf := bytes.NewBuffer(content)
	volumeInfo, err := VolumeInfoXMLUnmarshall(contentBuf)
	if err != nil {
		t.Errorf("Something went wrong while unmarshalling xml: %v", err)
	}

	if want, got := 0, volumeInfo.OpErrno; want != got {
		t.Errorf("want volumeInfo.OpErrno %d, got %d", want, got)
	}

	if want, got := 2, volumeInfo.VolInfo.Volumes.Count; want != got {
		t.Errorf("want volumeInfo.VolInfo.Volumes.Count %d, got %d", want, got)
	}

	volumeStrings := []string{"_all"}

	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if volumeStrings[0] == allVolumes || ContainsVolume(volumeStrings, volume.Name) {

			switch volume.Name {
			case "gv_cluster":
				if want, got := 4, volume.BrickCount; want != got {
					t.Errorf("want volume.BrickCount %d, got %d", want, got)
				}

				if want, got := 1, volume.Status; want != got {
					t.Errorf("want volume.Status %d, got %d", want, got)
				}
			case "gv_test":
				if want, got := 4, volume.BrickCount; want != got {
					t.Errorf("want volume.BrickCount %d, got %d", want, got)
				}

				if want, got := 1, volume.Status; want != got {
					t.Errorf("want volume.Status %d, got %d", want, got)
				}
			default: fmt.Printf("want %s or %s, got %s", "gv_cluster", "gv_test", "Error")
			}

		}
	}
}

func TestVolumeStatusAllDetailXMLUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/gluster_volume_status_all_detail.xml")
	if err != nil {
		t.Fatal(err)
	}

	// Convert into bytes.buffer
	contentBuf := bytes.NewBuffer(content)
	volumeStatusAll, err := VolumeStatusAllDetailXMLUnmarshall(contentBuf)
	if err != nil {
		t.Errorf("Something went wrong while unmarshalling xml: %v", err)
	}

	for _, vol := range volumeStatusAll.VolStatus.Volumes.Volume {

		for _, node := range vol.Node {
			if node.Status == 1 {

				if want, got := 20507914240, node.SizeTotal; want != int(got) {
					t.Errorf("want node.SizeTotal %d, got %d", want, got)
				}

				switch vol.VolName {
				case "gv_test":
					if want, got := "/mnt/gluster/gv_test", node.Path; want != got {
						t.Errorf("want node.Path %s, got %s", want, got)
					}
				case "gv_test2":
					if want, got := "/mnt/gluster/gv_test2", node.Path; want != got {
						t.Errorf("want node.Path %s, got %s", want, got)
					}
				default:
					t.Error("No vol.VolName match test instance")
				}

			}
		}
	}
}
