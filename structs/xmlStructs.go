package structs

import (
	"bytes"
	"encoding/xml"
	"github.com/prometheus/common/log"
	"io/ioutil"
)

type VolumeInfoXml struct {
	XMLName  xml.Name `xml:"cliOutput"`
	OpRet    int      `xml:"opRet"`
	OpErrno  int      `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	VolInfo  VolInfo  `xml:"volInfo"`
}

type VolInfo struct {
	XMLName xml.Name `xml:"volInfo"`
	Volumes Volumes  `xml:"volumes"`
}

type Volumes struct {
	XMLName xml.Name `xml:"volumes"`
	Volume  []Volume `xml:"volume"`
	Count   int      `xml:"count"`
}

type Volume struct {
	XMLName    xml.Name `xml:"volume"`
	Name       string   `xml:"name"`
	Id         string   `xml:"id"`
	Status     int      `xml:"status"`
	StatusStr  string   `xml:"statusStr"`
	BrickCount int      `xml:"brickCount"`
	Bricks     []Brick  `xml:"bricks"`
	DistCount  int      `xml:"distCount"`
}

type Brick struct {
	Uuid      string `xml:"brick>uuid"`
	Name      string `xml:"brick>name"`
	HostUuid  string `xml:"brick>hostUuid"`
	IsArbiter int    `xml:"brick>isArbiter"`
}

type VolumeListXml struct {
	XMLName  xml.Name `xml:"cliOutput"`
	OpRet    int      `xml:"opRet"`
	OpErrno  int      `xml:"opErrno"`
	OpErrstr string   `xml:"opErrstr"`
	VolList  VolList  `xml:"volList"`
}

type VolList struct {
	Count  int      `xml:"count"`
	Volume []string `xml:"volume"`
}

type PeerStatusXml struct {
	XMLName    xml.Name   `xml:"cliOutput"`
	OpRet      int        `xml:"opRet"`
	OpErrno    int        `xml:"opErrno"`
	OpErrstr   string     `xml:"opErrstr"`
	PeerStatus PeerStatus `xml:"peerStatus"`
}

type PeerStatus struct {
	XMLName xml.Name `xml:"peerStatus"`
	Peer    []Peer   `xml:"peer"`
}

type Peer struct {
	XMLName   xml.Name  `xml:"peer"`
	Uuid      string    `xml:"uuid"`
	Hostname  string    `xml:"hostname"`
	Hostnames Hostnames `xml:"hostnames"`
	Connected int       `xml:"connected"`
	State     int       `xml:"state"`
	StateStr  string    `xml:"stateStr"`
}

type Hostnames struct {
	Hostname string `xml:"hostname"`
}

type VolumeProfileXml struct {
	XMLName    xml.Name   `xml:"cliOutput"`
	OpRet      int        `xml:"opRet"`
	OpErrno    int        `xml:"opErrno"`
	OpErrstr   string     `xml:"opErrstr"`
	VolProfile VolProfile `xml:"volProfile"`
}

type VolProfile struct {
	Volname    string         `xml:"volname"`
	BrickCount int            `xml:"brickCount"`
	Brick      []BrickProfile `xml:"brick"`
}

type BrickProfile struct {
	//XMLName xml.Name `xml:"brick"`
	BrickName       string          `xml:"brickName"`
	CumulativeStats CumulativeStats `xml:"cumulativeStats"`
}

type CumulativeStats struct {
	//FopStats FopStats `xml:"fopStats"`
	Duration   int `xml:"duration"`
	TotalRead  int `xml:"totalRead"`
	TotalWrite int `xml:"totalWrite"`
}

type FopStats struct {
}

// Unmarshall returned bytes to VolumeListXml struct
func VolumeListXmlUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeListXml, error) {
	var vol VolumeListXml
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}

// Unmarshall returned bytes to VolumeInfoXml struct
func VolumeInfoXmlUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeInfoXml, error) {
	var vol VolumeInfoXml
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}

// Unmarshall returned bytes to PeerStatusXml struct
func PeerStatusXmlUnmarshall(cmdOutBuff *bytes.Buffer) (PeerStatusXml, error) {
	var vol PeerStatusXml
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}

// Unmarshall cumulative profile of gluster volume profile
func VolumeProfileGvInfoCumulativeXmlUnmarshall(cmdOutBuff *bytes.Buffer) (VolumeProfileXml, error) {
	var vol VolumeProfileXml
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	xml.Unmarshal(b, &vol)
	return vol, nil
}
