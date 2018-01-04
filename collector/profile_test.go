package collector

import (
	"testing"
	"io/ioutil"
	"bytes"
)

func TestVolumeProfileGvInfoCumulativeXMLUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/gluster_volume_profile_gv_test_info_cumulative.xml")
	if err != nil {
		t.Fatal(err)
	}

	// Convert into bytes.buffer
	contentBuf := bytes.NewBuffer(content)
	volumeProfile, err := VolumeProfileGvInfoCumulativeXMLUnmarshall(contentBuf)
	if err != nil {
		t.Errorf("Something went wrong while unmarshalling xml: %v", err)
	}

	for _, brick := range volumeProfile.VolProfile.Brick {

		switch brick.BrickName {
		// just test one node
		case "node1.example.local:/mnt/gluster/gv_test":
			if want, got := 16932, brick.CumulativeStats.Duration; want != got {
				t.Errorf("want brick.CumulativeStats.Duration %d, got %d", want, got)
			}

			if want, got := 0, brick.CumulativeStats.TotalRead; want != got {
				t.Errorf("want brick.CumulativeStats.TotalRead %d, got %d", want, got)
			}

			if want, got := 7590710, brick.CumulativeStats.TotalWrite; want != got {
				t.Errorf("want brick.CumulativeStats.TotalWrite %d, got %d", want, got)
			}

			for _, fop := range brick.CumulativeStats.FopStats.Fop {
				switch fop.Name {
				case "WRITE":
					if want, got := 58, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 224.500000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 183.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 807.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "STATFS":
					if want, got := 3, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 44.666667, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 32.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 69.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "FLUSH":
					if want, got := 1, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 117.000000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 117.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 117.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "GETXATTR":
					if want, got := 123, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 148.658537, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 17.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 1154.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "OPENDIR":
					if want, got := 87, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 4.091954, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 3.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 6.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "CREATE":
					if want, got := 1, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 23259.000000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 23259.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 23259.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "LOOKUP":
					if want, got := 119, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 68.495798, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 14.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 332.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "READDIR":
					if want, got := 174, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 1601.942529, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 195.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 4566.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "FINODELK":
					if want, got := 2, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 80.000000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 76.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 84.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "ENTRYLK":
					if want, got := 2, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 54.000000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 51.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 57.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				case "FXATTROP":
					if want, got := 2, fop.Hits; want != got {
						t.Errorf("want fop.Hits %d, got %d", want, got)
					}

					if want, got := 211.500000, fop.AvgLatency; want != got {
						t.Errorf("want fop.AvgLatency %d, got %d", want, got)
					}

					if want, got := 192.000000, fop.MinLatency; want != got {
						t.Errorf("want fop.MinLatency %d, got %d", want, got)
					}

					if want, got := 231.000000, fop.MaxLatency; want != got {
						t.Errorf("want fop.MaxLatency %d, got %d", want, got)
					}
				default:
					// just test one fop
					t.Error("No fop.Name match test instance")
				}

				break
			}
		default:
			// just test one node
			continue
		}
	}
}