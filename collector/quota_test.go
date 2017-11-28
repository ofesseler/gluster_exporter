package collector

import (
	"testing"
	"io/ioutil"
	"bytes"
)

func TestVolumeQuotaListXMLUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/gluster_volume_quota_list.xml")
	if err != nil {
		t.Fatal(err)
	}

	// Convert into bytes.buffer
	contentBuf := bytes.NewBuffer(content)
	volumeQuota, err := VolumeQuotaListXMLUnmarshall(contentBuf)
	if err != nil {
		t.Errorf("Something went wrong while unmarshalling xml: %v", err)
	}

	for _, limit := range volumeQuota.VolQuota.QuotaLimits {

		if want, got := 0.0, ExceededFunc(limit.SlExceeded); want != got {
			t.Errorf("want limit.SlExceeded %f, got %f", want, got)
		}

		if want, got := 0.0, ExceededFunc(limit.HlExceeded); want != got {
			t.Errorf("want limit.HlExceeded %f, got %f", want ,got)
		}

		switch limit.Path {
		case "/foo":
			if want, got := 10737418240, limit.HardLimit; want != int(got) {
				t.Errorf("want limit.HardLimit %d, got %d", want, got)
			}

			if want, got := 8589934592, limit.SoftLimitValue; want != int(got) {
				t.Errorf("want limit.SoftLimitValue %d, got %d", want, got)
			}

			if want, got := 428160000, limit.UsedSpace; want != int(got) {
				t.Errorf("want limit.UsedSpace %d, got %d", want, got)
			}

			if want, got := 10309258240, limit.AvailSpace; want != int(got) {
				t.Errorf("want limit.AvailSpace %d, got %d", want, got)
			}

		case "/bar":
			if want, got := 2147483648, limit.HardLimit; want != int(got) {
				t.Errorf("want limit.HardLimit %d, got %d", want, got)
			}

			if want, got := 1717986918, limit.SoftLimitValue; want != int(got) {
				t.Errorf("want limit.SoftLimitValue %d, got %d", want, got)
			}

			if want, got := 335544320, limit.UsedSpace; want != int(got) {
				t.Errorf("want limit.UsedSpace %d, got %d", want, got)
			}

			if want, got := 1811939328, limit.AvailSpace; want != int(got) {
				t.Errorf("want limit.AvailSpace %d, got %d", want, got)
			}
		default:
			t.Error("No limit.Path match test instance")
		}
	}
}
