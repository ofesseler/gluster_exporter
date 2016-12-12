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
