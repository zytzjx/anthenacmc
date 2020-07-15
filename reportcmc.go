package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"

	Log "github.com/zytzjx/anthenacmc/loggersys"
)

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	// return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func transcation(url string, info map[string]interface{}) {
	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(info).
		Post(fmt.Sprintf("%s%s", url, "api/auth/"))

	if err != nil {
		fmt.Println("web request fail")
		Log.Log.Error("web request fail")
		return
	}
	Log.Log.Info(string(resp.Body()))

	var status map[string]interface{}
	fmt.Println(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), &status); err != nil {
		fmt.Println("return format error.")
		Log.Log.Error("return format error.")
	}

	if val, ok := status["status"]; ok {
		switch val {
		case "1", "2":
			fmt.Println("success")
			Log.Log.Info("Success")
		case "4":
			Log.Log.Warn("same uuid")
			fmt.Println("same uuid")
		default:
			fmt.Println("failed.")
			Log.Log.Error("failed.")
		}
	}

}
