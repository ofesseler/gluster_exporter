package collector

import (
	"bytes"
	"io/ioutil"
	"encoding/xml"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// Subsystem(s).
	peer = "peer"
)

var (
	peersConnected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, peer, "peers_connected"),
		"Is peer connected to gluster cluster.",
		nil, nil,
	)
	peersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, peer, "peers_total"),
		"Total peer nums of gluster cluster.",
		nil, nil,
	)
)

func ScrapePeerStatus(ch chan<- prometheus.Metric) error {
	// Read gluster peer status
	peerStatus, peerStatusErr := ExecPeerStatus()
	if peerStatusErr != nil {
		log.Errorf("Couldn't parse xml of peer status: %v", peerStatusErr)
		return peerStatusErr
	}

	countConnected := 0
	countTotal := 0
	for _, peer := range peerStatus.Peer {
		// State 3 means "Peer in Cluster"
		if peer.Connected == 1 && peer.State == 3 {
			countConnected++
		}
		countTotal++
	}

	ch <- prometheus.MustNewConstMetric(
		peersConnected, prometheus.GaugeValue, float64(countConnected),
	)

	ch <- prometheus.MustNewConstMetric(
		peersTotal, prometheus.GaugeValue, float64(countTotal),
	)

	return nil
}

// ExecPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func ExecPeerStatus() (PeerStatus, error) {
	args := []string{"peer", "status"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return PeerStatus{}, cmdErr
	}
	peerStatus, err := PeerStatusXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

// PeerStatusXMLUnmarshall unmarshalls bytes to PeerStatusXML struct
func PeerStatusXMLUnmarshall(cmdOutBuff *bytes.Buffer) (PeerStatusXML, error) {
	var vol PeerStatusXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}
