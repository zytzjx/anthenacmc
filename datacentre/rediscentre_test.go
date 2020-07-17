package datacentre

import "testing"

func TestGetSerialConfig(t *testing.T) {
	if _, err := GetSerialConfig(); err != nil {
		t.Error(err)
	}
	t.Log("Success")
}
