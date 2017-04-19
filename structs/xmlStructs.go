package structs

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"

	"github.com/prometheus/common/log"
)

// VolumeInfoXML struct repesents cliOutput element of "gluster volume info" command
type VolumeInfoXML struct {
	XMLName  xml.Name `xml:"cliOutput"`
	OpRet    int      `xml:"opRet"`
	OpErrno  int      `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	VolInfo  VolInfo  `xml:"volInfo"`
}

// VolInfo element of "gluster volume info" command
type VolInfo struct {
	XMLName xml.Name `xml:"volInfo"`
	Volumes Volumes  `xml:"volumes"`
}

// Volumes element of "gluster volume info" command
type Volumes struct {
	XMLName xml.Name `xml:"volumes"`
	Volume  []Volume `xml:"volume"`
	Count   int      `xml:"count"`
}

// Volume element of "gluster volume info" command
type Volume struct {
	XMLName    xml.Name `xml:"volume"`
	Name       string   `xml:"name"`
	ID         string   `xml:"id"`
	Status     int      `xml:"status"`
	StatusStr  string   `xml:"statusStr"`
	BrickCount int      `xml:"brickCount"`
	Bricks     []Brick  `xml:"bricks"`
	DistCount  int      `xml:"distCount"`
}

// Brick element of "gluster volume info" command
type Brick struct {
	UUID      string `xml:"brick>uuid"`
	Name      string `xml:"brick>name"`
	HostUUID  string `xml:"brick>hostUuid"`
	IsArbiter int    `xml:"brick>isArbiter"`
}

// VolumeListXML struct repesents cliOutput element of "gluster volume list" command
type VolumeListXML struct {
	XMLName  xml.Name `xml:"cliOutput"`
	OpRet    int      `xml:"opRet"`
	OpErrno  int      `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	VolList  VolList  `xml:"volList"`
}

// VolList element of "gluster volume list" command
type VolList struct {
	Count  int      `xml:"count"`
	Volume []string `xml:"volume"`
}

// PeerStatusXML struct represents cliOutput element of "gluster peer status" command
type PeerStatusXML struct {
	XMLName    xml.Name   `xml:"cliOutput"`
	OpRet      int        `xml:"opRet"`
	OpErrno    int        `xml:"opErrno"`
	OpErrstr   string     `xml:"opErrstr"`
	PeerStatus PeerStatus `xml:"peerStatus"`
}

// PeerStatus element of "gluster peer status" command
type PeerStatus struct {
	XMLName xml.Name `xml:"peerStatus"`
	Peer    []Peer   `xml:"peer"`
}

// Peer element of "gluster peer status" command
type Peer struct {
	XMLName   xml.Name  `xml:"peer"`
	UUID      string    `xml:"uuid"`
	Hostname  string    `xml:"hostname"`
	Hostnames Hostnames `xml:"hostnames"`
	Connected int       `xml:"connected"`
	State     int       `xml:"state"`
	StateStr  string    `xml:"stateStr"`
}

// Hostnames element of "gluster peer status" command
type Hostnames struct {
	Hostname string `xml:"hostname"`
}

// VolumeProfileXML struct repesents cliOutput element of "gluster volume {volume} profile" command
type VolumeProfileXML struct {
	XMLName    xml.Name   `xml:"cliOutput"`
	OpRet      int        `xml:"opRet"`
	OpErrno    int        `xml:"opErrno"`
	OpErrstr   string     `xml:"opErrstr"`
	VolProfile VolProfile `xml:"volProfile"`
}

// VolProfile element of "gluster volume {volume} profile" command
type VolProfile struct {
	Volname    string         `xml:"volname"`
	BrickCount int            `xml:"brickCount"`
	Brick      []BrickProfile `xml:"brick"`
}

// BrickProfile struct for element brick of "gluster volume {volume} profile" command
type BrickProfile struct {
	//XMLName xml.Name `xml:"brick"`
	BrickName       string          `xml:"brickName"`
	CumulativeStats CumulativeStats `xml:"cumulativeStats"`
}

// CumulativeStats element of "gluster volume {volume} profile" command
type CumulativeStats struct {
	FopStats   FopStats `xml:"fopStats"`
	Duration   int      `xml:"duration"`
	TotalRead  int      `xml:"totalRead"`
	TotalWrite int      `xml:"totalWrite"`
}

// FopStats element of "gluster volume {volume} profile" command
type FopStats struct {
	Fop []Fop `xml:"fop"`
}

// Fop is struct for FopStats
type Fop struct {
	Name       string  `xml:"name"`
	Hits       int     `xml:"hits"`
	AvgLatency float64 `xml:"avgLatency"`
	MinLatency float64 `xml:"minLatency"`
	MaxLatency float64 `xml:"maxLatency"`
}

type HealInfoBrick struct {
	XMLName         xml.Name `xml:"brick"`
	Name            string   `xml:"name"`
	Status          string   `xml:"status"`
	NumberOfEntries string   `xml:"numberOfEntries"`
}

type HealInfoBricks struct {
	XMLName xml.Name        `xml:"bricks"`
	Brick   []HealInfoBrick `xml:"brick"`
}

type HealInfo struct {
	XMLName xml.Name       `xml:"healInfo"`
	Bricks  HealInfoBricks `xml:"bricks"`
}

// VolumeHealInfoXML struct repesents cliOutput element of "gluster volume {volume} heal info" command
type VolumeHealInfoXML struct {
	XMLName  xml.Name `xml:"cliOutput"`
	OpRet    int      `xml:"opRet"`
	OpErrno  int      `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	HealInfo HealInfo `xml:"healInfo"`
}

// VolumeHealInfoXMLUnmarshall unmarshalls heal info of gluster cluster
func VolumeHealInfoXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeHealInfoXML, error) {
	var vol VolumeHealInfoXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	err = xml.Unmarshal(b, &vol)
	if err != nil {
		log.Error(err)
	}
	return vol, nil
}

// VolumeListXMLUnmarshall unmarshalls bytes to VolumeListXML struct
func VolumeListXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeListXML, error) {
	var vol VolumeListXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
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

// VolumeProfileGvInfoCumulativeXMLUnmarshall unmarshalls cumulative profile of gluster volume profile
func VolumeProfileGvInfoCumulativeXMLUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeProfileXML, error) {
	var vol VolumeProfileXML
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}

// VolumeStatusXML XML type of "gluster volume status"
type VolumeStatusXML struct {
	XMLName   xml.Name `xml:"cliOutput"`
	OpRet     int      `xml:"opRet"`
	OpErrno   int      `xml:"opErrno"`
	OpErrstr  string   `xml:"opErrstr"`
	VolStatus struct {
		Volumes []struct {
			Volume struct {
				VolName   string `xml:"volName"`
				NodeCount int    `xml:"nodeCount"`
				Node      []struct {
					Hostname string `xml:"hostname"`
					Path     string `xml:"path"`
					PeerID   string `xml:"peerid"`
					Status   int    `xml:"status"`
					Port     int    `xml:"port"`
					Ports    struct {
						TCP  int    `xml:"tcp"`
						RDMA string `xml:"rdma"`
					} `xml:"ports"`
					Pid        int    `xml:"pid"`
					SizeTotal  uint64 `xml:"sizeTotal"`
					SizeFree   uint64 `xml:"sizeFree"`
					Device     string `xml:"device"`
					BlockSize  int    `xml:"blockSize"`
					MntOptions string `xml:"mntOptions"`
					FsName     string `xml:"fsName"`
				} `xml:"node"`
			} `xml:"volume"`
		} `xml:"volumes"`
	} `xml:"volStatus"`
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

type QuotaLimit struct {
    XMLName          xml.Name `xml:"limit"`
    Path             string   `xml:"path"`
    HardLimit        uint64   `xml:"hard_limit"`
    SoftLimitValue   uint64   `xml:"soft_limit_value"`
    UsedSpace        uint64   `xml:"used_space"`
    AvailSpace       uint64   `xml:"avail_space"`
    SlExceeded       string   `xml:"sl_exceeded"`
    HlExceeded       string   `xml:"hl_exceeded"`
}

type VolQuota struct {
    XMLName     xml.Name     `xml:"volQuota"`
    QuotaLimits []QuotaLimit `xml:"limit"`
}
// VolumeQuotaXML XML type of "gluster volume quota list"
type VolumeQuotaXML struct {
    XMLName  xml.Name  `xml:"cliOutput"`
    OpRet     int      `xml:"opRet"`
    OpErrno   int      `xml:"opErrno"`
    OpErrstr  string   `xml:"opErrstr"`
    VolQuota  VolQuota `xml:"volQuota"`
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
