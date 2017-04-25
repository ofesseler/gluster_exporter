package structs

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
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
	t.Log("Test xml unmarshal for 'gluster volume status all detail' with file: ", testXMLPath)
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

	for _, vol := range volumeStatus.VolStatus.Volumes.Volume {
		if vol.NodeCount != 4 {
			t.Errorf("nodecount mismatch %v instead of 4", vol.NodeCount)
		}

		for _, node := range vol.Node {
			if node.BlockSize != 4096 {
				t.Errorf("blockSize mismatch %v and 4096 expected", node.BlockSize)
			}

		}

		if vol.Node[0].SizeFree != 19517558784 {
			t.Errorf("SizeFree doesn't match 19517558784: %v", vol.Node[0].SizeFree)
		}

		if vol.Node[0].SizeTotal != 20507914240 {
			t.Errorf("SizeFree doesn't match 20507914240: %v", vol.Node[0].SizeTotal)
		}
	}

	if volumeStatus.VolStatus.Volumes.Volume[0].VolName != "gv_test" {
		t.Errorf("VolName of first volume doesn't match gv_test: %v", volumeStatus.VolStatus.Volumes.Volume[0].VolName)
	}

	if volumeStatus.VolStatus.Volumes.Volume[1].VolName != "gv_test2" {
		t.Errorf("VolName of first volume doesn't match gv_test2: %v", volumeStatus.VolStatus.Volumes.Volume[1].VolName)
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

	if fops[0].AvgLatency != expAvgLatency {
		t.Errorf("expected %v as name and got %v", expAvgLatency, fops[0].AvgLatency)
	}

	if fops[0].MinLatency != expMinLatency {
		t.Errorf("expected %v as name and got %v", expMinLatency, fops[0].MinLatency)
	}

	if fops[0].MaxLatency != expMaxLatency {
		t.Errorf("expected %v as name and got %v", expMaxLatency, fops[0].MaxLatency)
	}
}

func getCliBufferHelper(filename string) *bytes.Buffer {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Could not read test data from xml.", err)
	}
	return bytes.NewBuffer(dat)
}

type testPair struct {
	path      string
	expected  int
	nodeCount int
}

func TestVolumeHealInfoXMLUnmarshall(t *testing.T) {
	var test = []testPair{
		{path: "../test/gluster_volume_heal_info_err_node2.xml", expected: 3, nodeCount: 4},
	}

	for _, c := range test {
		cmdOutBuffer := getCliBufferHelper(c.path)
		healInfo, err := VolumeHealInfoXMLUnmarshall(cmdOutBuffer)
		if err != nil {
			t.Error(err)
		}
		if healInfo.OpErrno != 0 {
			t.Error(healInfo.OpErrstr)
		}
		entriesOutOfSync := 0
		if len(healInfo.HealInfo.Bricks.Brick) != c.nodeCount {
			t.Error(healInfo.HealInfo.Bricks)
			t.Errorf("Excpected %v Bricks and len is %v", c.nodeCount, len(healInfo.HealInfo.Bricks.Brick))
		}
		for _, brick := range healInfo.HealInfo.Bricks.Brick {
			var count int
			count, _ = strconv.Atoi(brick.NumberOfEntries)
			entriesOutOfSync += count
		}
		if entriesOutOfSync != c.expected {
			t.Errorf("Out of sync entries other than expected: %v and was %v", c.expected, entriesOutOfSync)
		}
	}

}

func TestVolumeQuotaListXMLUnmarshall(t *testing.T) {
    testXMLPath := "../test/gluster_volume_quota_list.xml"
    nodeCount := 2
    dat, err := ioutil.ReadFile(testXMLPath)

	if err != nil {
		t.Errorf("error reading testxml in Path: %v", testXMLPath)
	}
	volumeQuotaXML, err := VolumeQuotaListXMLUnmarshall(bytes.NewBuffer(dat))
	if err != nil {
		t.Error(err)
	}

	if volumeQuotaXML.OpErrno != 0 {
		t.Error(volumeQuotaXML.OpErrstr)
	}
    nb_limits := len(volumeQuotaXML.VolQuota.QuotaLimits)
    if nb_limits != nodeCount {
        t.Errorf("Expected %v Limits and len is %v", nodeCount, nb_limits)
    }

    for _, limit := range volumeQuotaXML.VolQuota.QuotaLimits {
        if limit.Path == "/foo" {
            if limit.AvailSpace != 10309258240 {
                t.Errorf(
                    "Expected %v for available space in path %v, got %v",
                    1811939328,
                    limit.Path,
                    limit.AvailSpace,
                )
            }
        }
    }

}
