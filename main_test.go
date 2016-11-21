package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestInfoUnmarshall(t *testing.T) {
	testXmlPath := "test/gluster_volume_info.xml"
	t.Log("Test xml unmarshal for gluster volume info with file: ", testXmlPath)
	dat, err := ioutil.ReadFile(testXmlPath)

	if err != nil {
		t.Errorf("error reading testxml in Path: %v", testXmlPath)
	}

	glusterVolumeInfo := infoUnmarshall(bytes.NewBuffer(dat))
	if glusterVolumeInfo.OpErrno != 0 && glusterVolumeInfo.VolInfo.Volumes.Count == 2 {
		t.Error("something wrong")
	}
	t.Log("gluster volume info test was sucessful.")
}
