package collector

import "encoding/xml"

// VolumeInfoXML struct repesents cliOutput element of "gluster volume info" command
//
// cliOutput
// |-- opRet
// |-- opErrno
// |-- opErrstr
// |-- volInfo
//     |-- volumes
//         |-- volume
//             |-- name
// 			   |-- id
// 			   |-- status
// 			   |-- statusStr
// 			   |-- snapshotCount
// 			   |-- brickCount
// 			   |-- distCount
// 			   |-- stripeCount
// 			   |-- replicaCount
// 			   |-- arbiterCount
// 			   |-- disperseCount
// 			   |-- redundancyCount
// 			   |-- type
// 			   |-- typeStr
// 			   |-- transport
//  		   |-- xlators/       // TODO: don't know what means
// 			   |-- bricks
// 			       |-- []brick
// 			           |-- name
// 			           |-- hostUuid
// 			           |-- isArbiter
// 			   |-- optCount
//             |-- options
//                 |-- []option
//                     |-- name
//                     |-- value

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


// PeerStatusXML struct represents cliOutput element of "gluster peer status" command
//
// cliOutput
// |-- opRet
// |-- opErrno
// |-- opErrstr
// |-- peerStatus
//     |-- []peer
//         |-- uuid
//         |-- hostname
//         |-- hostnames
// 		       |-- hostname
//         |-- connected
//         |-- state
//         |-- stateStr

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

// VolumeStatusXML XML type of "gluster volume status"
//
// cliOutput
// |-- opRet
// |-- opErrno
// |-- opErrstr
// |-- volStatus
//     |-- volumes
//         |-- []volume
//         	   |-- volName
//         	   |-- nodeCount
//         	   |-- []node
// 			       |-- hostname
// 			       |-- path
// 			       |-- peerid
// 			       |-- status
// 			       |-- port
// 			       |-- ports
// 			       	   |-- tcp
// 			       	   |-- rdma
// 			       |-- pid

type VolumeStatusXML struct {
	XMLName   xml.Name `xml:"cliOutput"`
	OpRet     int      `xml:"opRet"`
	OpErrno   int      `xml:"opErrno"`
	OpErrstr  string   `xml:"opErrstr"`
	VolStatus struct {
		Volumes struct {
			Volume []struct {
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

// Quota
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
// VolumeQuotaXML XML type of "gluster volume quota {volume} list"
type VolumeQuotaXML struct {
	XMLName  xml.Name  `xml:"cliOutput"`
	OpRet     int      `xml:"opRet"`
	OpErrno   int      `xml:"opErrno"`
	OpErrstr  string   `xml:"opErrstr"`
	VolQuota  VolQuota `xml:"volQuota"`
}

// Profile
// VolumeProfileXML struct repesents cliOutput element of "gluster volume profile {volume} info" command
//
// cliOutput
// |-- opRet
// |-- opErrno
// |-- opErrstr
// |-- volProfile
//     |-- volname
//     |-- profileOp
//     |-- brickCount
//     |-- []brick
//         |-- brickName
//         |-- cumulativeStats
//         	   |-- blockStats
// 		   	      |-- []block
//             |-- fopStats
// 			      |-- []fop
//             |-- duration
//             |-- totalRead
//             |-- totalWrite

type VolumeProfileXML struct {
	XMLName    xml.Name   `xml:"cliOutput"`
	OpRet      int        `xml:"opRet"`
	OpErrno    int        `xml:"opErrno"`
	OpErrstr   string     `xml:"opErrstr"`
	VolProfile VolProfile `xml:"volProfile"`
}

// VolProfile element of "gluster volume profile {volume} info" command
type VolProfile struct {
	Volname    string         `xml:"volname"`
	ProfileOp  int 			  `xml:"profileOp"`
	BrickCount int            `xml:"brickCount"`
	Brick      []BrickProfile `xml:"brick"`
}

// BrickProfile struct for element brick of "gluster volume profile {volume} info" command
type BrickProfile struct {
	//XMLName xml.Name `xml:"brick"`
	BrickName       string          `xml:"brickName"`
	CumulativeStats CumulativeStats `xml:"cumulativeStats"`
}

// CumulativeStats element of "gluster volume profile {volume} info" command
type CumulativeStats struct {
	FopStats   FopStats `xml:"fopStats"`
	Duration   int      `xml:"duration"`
	TotalRead  int      `xml:"totalRead"`
	TotalWrite int      `xml:"totalWrite"`
}

// FopStats element of "gluster volume profile {volume} info" command
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