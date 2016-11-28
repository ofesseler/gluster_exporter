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
