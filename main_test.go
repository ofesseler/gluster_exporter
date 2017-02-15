package main

import "testing"

func TestContainsVolume(t *testing.T) {
	expamle := "doge"
	testSlice := []string{"wow", "such", expamle}
	if !ContainsVolume(testSlice, expamle) {
		t.Fatalf("Hasn't found %v in slice %v", expamle, testSlice)
	}
}

func TestParseMountOutput(t *testing.T) {
	mountOutput := "/dev/mapper/cryptroot on / type ext4 (rw,relatime,data=ordered) \n /dev/mapper/cryptroot on /var/lib/docker/devicemapper type ext4 (rw,relatime,data=ordered)"
	mounts, err := parseMountOutput("asd", mountOutput)
	if err != nil {
		t.Error(err)
	}
	expected := []string{"/", "/var/lib/docker/devicemapper"}
	for i, mount := range mounts {
		if mount.mountPoint != expected[i] {
			t.Errorf("mountpoint is %v and %v was expected", mount.mountPoint, expected[i])
		}
	}
}
