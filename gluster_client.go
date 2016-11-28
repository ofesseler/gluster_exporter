package main

import (
	"bytes"
	"github.com/ofesseler/gluster_exporter/structs"
	"github.com/prometheus/common/log"
	"os/exec"
)

func execGlusterCommand(arg ...string) *bytes.Buffer {
	stdoutBuffer := &bytes.Buffer{}
	argXML := append(arg, "--xml")
	glusterExec := exec.Command(GlusterCmd, argXML...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		log.Fatal(err)
	}
	return stdoutBuffer
}

// ExecVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeInfo() (structs.VolumeInfoXML, error) {
	args := []string{"volume", "info"}
	bytesBuffer := execGlusterCommand(args...)
	volumeInfo, err := structs.VolumeInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// ExecVolumeList executes "gluster volume info" at the local machine and
// returns VolumeList struct and error
func ExecVolumeList() (structs.VolList, error) {
	args := []string{"volume", "list"}
	bytesBuffer := execGlusterCommand(args...)
	volumeList, err := structs.VolumeListXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeList.VolList, err
	}

	return volumeList.VolList, nil
}

// ExecPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func ExecPeerStatus() (structs.PeerStatus, error) {
	args := []string{"peer", "status"}
	bytesBuffer := execGlusterCommand(args...)
	peerStatus, err := structs.PeerStatusXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

// ExecVolumeProfileGvInfoCumulative executes "gluster volume {volume] profile info cumulative" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeProfileGvInfoCumulative(volumeName string) (structs.VolProfile, error) {
	args := []string{"volume", "profile"}
	args = append(args, volumeName)
	args = append(args, "info", "cumulative")
	bytesBuffer := execGlusterCommand(args...)
	volumeProfile, err := structs.VolumeProfileGvInfoCumulativeXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}
