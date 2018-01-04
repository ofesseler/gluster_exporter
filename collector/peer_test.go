package collector

import (
	"testing"
	"io/ioutil"
	"bytes"
)

func TestPeerStatusXMLUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/gluster_peer_status.xml")
	if err != nil {
		t.Fatal(err)
	}

	// Convert into bytes.buffer
	contentBuf := bytes.NewBuffer(content)
	peerStatus, err := PeerStatusXMLUnmarshall(contentBuf)
	if err != nil {
		t.Errorf("Something went wrong while unmarshalling xml: %v", err)
	}

	count := 0
	for _, peer := range peerStatus.PeerStatus.Peer {
		if peer.Connected == 1 && peer.State == 3 {
			count ++
		}
	}

	if want, got := 3, count; want != got {
		t.Errorf("want peer count %d, got %d", want, got)
	}

}