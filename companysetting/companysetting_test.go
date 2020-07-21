package companysetting

import "testing"

func TestGetRemoteIP(t *testing.T) {
	ip, err := getRemoteIP()
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}

func TestGetLocalPCInfo(t *testing.T) {
	mac, ip, err := GetLocalPCInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(mac, ip)
}

func TestLoadIPInfoFromFile(t *testing.T) {
	mac, ip, err := loadIPInfoFromFile()
	if err != nil {
		t.Error(err)
	}
	t.Log(mac, ip)
}

func TestDownload(t *testing.T) {
	Download()
}
