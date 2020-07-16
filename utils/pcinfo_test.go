package utils

import (
	"fmt"
	"testing"
)

func TestGetMacAddr(t *testing.T) {
	addrs, err := GetMacAddr()
	if err != nil {
		t.Error("get mac addr error")
	}
	if len(addrs) > 0 {
		t.Log("Success")
	} else {
		t.Fail()
	}
}

func TestGetMapMacIP(t *testing.T) {
	macips, err := GetMapMacIP()
	if err != nil {
		t.Error("get mac addr error")
	}
	for k, v := range macips {
		fmt.Println(k, "     ", v)
	}
}

func TestGetIPAddrs(t *testing.T) {
	ips, err := GetIPAddrs()
	if err != nil {
		t.Error("get mac addr error")
	}
	if len(ips) > 0 {
		for _, ip := range ips {
			fmt.Println(ip)
		}
		t.Log("Success")
	} else {
		t.Fail()
	}
}

func TestGetPCName(t *testing.T) {
	GetPCName()
}
