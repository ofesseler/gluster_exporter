package main

import "testing"

func TestContainsVolume(t *testing.T) {
	expamle := "doge"
	testSlice := []string{"wow", "such", expamle}
	if !ContainsVolume(testSlice, expamle) {
		t.Fatalf("Hasn't found %v in slice %v", expamle, testSlice)
	}
}

type testCases struct {
	mountOutput string
	expected    []string
}

func TestParseMountOutput(t *testing.T) {
	var tests = []testCases{
		{
			mountOutput: "/dev/mapper/cryptroot on / type ext4 (rw,relatime,data=ordered) \n" +
				"/dev/mapper/cryptroot on /var/lib/docker/devicemapper type ext4 (rw,relatime,data=ordered)",
			expected: []string{"/", "/var/lib/docker/devicemapper"},
		},
		{
			mountOutput: "/dev/mapper/cryptroot on / type ext4 (rw,relatime,data=ordered) \n" +
				"",
			expected: []string{"/"},
		},
	}
	for _, c := range tests {
		mounts, err := parseMountOutput(c.mountOutput)
		if err != nil {
			t.Error(err)
		}

		for i, mount := range mounts {
			if mount.mountPoint != c.expected[i] {
				t.Errorf("mountpoint is %v and %v was expected", mount.mountPoint, c.expected[i])
			}
		}
	}

}
