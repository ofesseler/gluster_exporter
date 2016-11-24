package main

import (
	"bytes"
	"github.com/ofesseler/gluster_exporter/structs"
	"github.com/prometheus/common/log"
	"os/exec"
)

func execGlusterCommand(arg ...string) *bytes.Buffer {
	stdoutBuffer := &bytes.Buffer{}
	arg_xml := append(arg, "--xml")
	glusterExec := exec.Command(GLUSTER_CMD, arg_xml...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		log.Fatal(err)
	}
	return stdoutBuffer
}

func ExecVolumeInfo() (structs.VolumeInfoXml, error) {
	args := []string{"volume", "info"}
	bytesBuffer := execGlusterCommand(args...)
	volumeInfo, err := structs.VolumeInfoXmlUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

func ExecVolumeList() (structs.VolList, error) {
	args := []string{"volume", "list"}
	bytesBuffer := execGlusterCommand(args...)
	volumeList, err := structs.VolumeListXmlUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeList.VolList, err
	}

	return volumeList.VolList, nil
}

func ExecPeerStatus() (structs.PeerStatus, error) {
	args := []string{"peer", "status"}
	bytesBuffer := execGlusterCommand(args...)
	peerStatus, err := structs.PeerStatusXmlUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

func ExecVolumeProfileGvInfoCumulative(volumeName string) (structs.VolProfile, error) {
	args := []string{"volume", "profile"}
	args = append(args, volumeName)
	args = append(args, "info", "cumulative")
	bytesBuffer := execGlusterCommand(args...)
	volumeProfile, err := structs.VolumeProfileGvInfoCumulativeXmlUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}
