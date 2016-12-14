package structs

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestVolumeListXMLUnmarshall(t *testing.T) {
	testXMLPath := "../test/gluster_volume_list.xml"
	dat, err := ioutil.ReadFile(testXMLPath)

	if err != nil {
		t.Errorf("error reading testxml in Path: %v", testXMLPath)
	}
	volumeList, err := VolumeListXMLUnmarshall(bytes.NewBuffer(dat))
	if err != nil {
		t.Fatal(err)
	}

	if len(volumeList.VolList.Volume) != 2 {
		t.Fatal("Volume List empty")
	}
	if volumeList.VolList.Count != 2 {
		t.Logf("doesn't match volume count of 2: %v", volumeList)
	}

	t.Log("gluster volume list test was successful.")
}

func TestInfoUnmarshall(t *testing.T) {
	testXMLPath := "../test/gluster_volume_info.xml"
	dat, err := ioutil.ReadFile(testXMLPath)

	if err != nil {
		t.Fatalf("error reading testxml in Path: %v", testXMLPath)
	}

	glusterVolumeInfo, _ := VolumeInfoXMLUnmarshall(bytes.NewBuffer(dat))
	if glusterVolumeInfo.OpErrno != 0 && glusterVolumeInfo.VolInfo.Volumes.Count == 2 {
		t.Fatal("something wrong")
	}
	for _, volume := range glusterVolumeInfo.VolInfo.Volumes.Volume {
		if volume.Status != 1 {
			t.Errorf("Status %v expected but got %v", 1, volume.Status)
		}
		if volume.Name == "" || len(volume.Name) < 1 {
			t.Errorf("Not empty name of Volume expected, response was %v", volume.Name)
		}
		t.Logf("Volume.Name: %v volume.Status: %v", volume.Name, volume.Status)
	}
	t.Log("gluster volume info test was successful.")
}

func TestPeerStatusXMLUnmarshall(t *testing.T) {
	testXMLPath := "../test/gluster_peer_status.xml"
	t.Log("Test xml unmarshal for 'gluster peer status' with file: ", testXMLPath)
	dat, err := ioutil.ReadFile(testXMLPath)
	exp := 0
	if err != nil {
		t.Errorf("error reading testxml in Path: %v", testXMLPath)
	}
	peerStatus, err := PeerStatusXMLUnmarshall(bytes.NewBuffer(dat))
	if err != nil {
		t.Fatal(err)
	}

	exp = 0
	if peerStatus.OpErrno != exp {
		t.Fatalf("OpErrno: %v", peerStatus.OpErrno)
	}

	exp = 1
	if peerStatus.PeerStatus.Peer[0].Connected != exp {
		t.Fatalf("Peerstatus is not as expected. Expectel: %v, Current: %v", exp, peerStatus.PeerStatus.Peer[0].Connected)
	}

	exp = 3
	if len(peerStatus.PeerStatus.Peer) != exp {
		t.Fatalf("Number of peers is not 3: %v", len(peerStatus.PeerStatus.Peer))
	}

	expString := "node2.example.local"
	if peerStatus.PeerStatus.Peer[0].Hostname != expString {
		t.Fatalf("Hostname in Peer does't match: %v", peerStatus.PeerStatus.Peer[0].Hostname)
	}

	t.Log("gluster peer status test was successful.")
}

func TestVolumeStatusAllDetailXMLUnmarshall(t *testing.T) {
	testXMLPath := "../test/gluster_volume_status_all_detail.xml"
	t.Log("Test xml unmarshal for 'gluster peer status' with file: ", testXMLPath)
	dat, err := ioutil.ReadFile(testXMLPath)
	if err != nil {
		t.Errorf("error reading testxml in Path: %v", testXMLPath)
	}
	volumeStatus, err := VolumeStatusAllDetailXMLUnmarshall(bytes.NewBuffer(dat))
	if err != nil {
		t.Error(err)
	}

	if volumeStatus.OpErrno != 0 {
		t.Error(volumeStatus.OpErrstr)
	}

	for _, vol := range volumeStatus.VolStatus.Volumes {
		if vol.Volume.NodeCount != 4 {
			t.Errorf("nodecount mismatch %v instead of 4", vol.Volume.NodeCount)
		}

		for _, node := range vol.Volume.Node {
			if node.BlockSize != 4096 {
				t.Errorf("blockSize mismatch %v and 4096 expected", node.BlockSize)
			}

		}

		if vol.Volume.Node[0].SizeFree != 19517558784 {
			t.Errorf("SizeFree doesn't match 19517558784: %v", vol.Volume.Node[0].SizeFree)
		}

		if vol.Volume.Node[0].SizeTotal != 20507914240 {
			t.Errorf("SizeFree doesn't match 20507914240: %v", vol.Volume.Node[0].SizeTotal)
		}
	}
}

func TestVolumeProfileGvInfoCumulativeXMLUnmarshall(t *testing.T) {
	testXMLPath := "../test/gluster_volume_profile_gv_test_info_cumulative.xml"
	t.Log("Test xml unmarshal for 'gluster volume profile gv_test info' with file: ", testXMLPath)
	dat, err := ioutil.ReadFile(testXMLPath)

	if err != nil {
		t.Fatal("Could not read test data from xml.", err)
	}

	profileVolumeCumulative, err := VolumeProfileGvInfoCumulativeXMLUnmarshall(bytes.NewBuffer(dat))
	if err != nil {
		t.Fatal(err)
	}

	expOpErr := 0
	if profileVolumeCumulative.OpErrno != 0 {
		t.Errorf("Expected value is %v and got %v", expOpErr, profileVolumeCumulative.OpErrno)
	}

	fops := profileVolumeCumulative.VolProfile.Brick[0].CumulativeStats.FopStats.Fop

	expFopLen := 11
	fopLen := len(fops)
	if fopLen != expFopLen {
		t.Errorf("Expected FopLength of %v and got %v.", expFopLen, fopLen)
	}

	expFopName := "WRITE"
	expFopHits := 58
	expAvgLatency := 224.5
	expMinLatency := 183.0
	expMaxLatency := 807.0

	if fops[0].Name != expFopName {
		t.Errorf("expected %v as name and got %v", expFopName, fops[0].Name)
	}

	if fops[0].Hits != expFopHits {
		t.Errorf("expected %v as name and got %v", expFopHits, fops[0].Hits)
	}

	if fops[0].AvgLatency!= expAvgLatency {
		t.Errorf("expected %v as name and got %v", expAvgLatency, fops[0].AvgLatency)
	}

	if fops[0].MinLatency != expMinLatency {
		t.Errorf("expected %v as name and got %v", expMinLatency, fops[0].MinLatency)
	}

	if fops[0].MaxLatency != expMaxLatency {
		t.Errorf("expected %v as name and got %v", expMaxLatency, fops[0].MaxLatency)
	}
}
