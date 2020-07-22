package companysetting

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	dmc "github.com/zytzjx/anthenacmc/datacentre"
	Log "github.com/zytzjx/anthenacmc/loggersys"
	Util "github.com/zytzjx/anthenacmc/utils"
)

type _Status struct{}
type _Sync struct {
	Status _Status `json:"status"`
}

type _Extrainfo struct {
	IP         string `json:"ip"`
	Product    string `json:"product"`
	ID         string `json:"_id"`
	RemoteAddr string `json:"remote_addr"`
}

// ClientInfo get company setting
type ClientInfo struct {
	Company       string     `json:"company"`
	Solutionid    string     `json:"solutionid"`
	Productid     string     `json:"productid"`
	Site          string     `json:"site"`
	ID            string     `json:"_id"`
	PCName        string     `json:"pcname"`
	Macaddr       string     `json:"macaddr"`
	Filecorrupted string     `json:"filecorrupted"`
	Extrainfo     _Extrainfo `json:"extrainfo"`
}

// RequestSetting get company setting request body
type RequestSetting struct {
	Client   ClientInfo `json:"client"`
	Sync     _Sync      `json:"sync"`
	Protocol string     `json:"protocol"`
}

// LocalMacIP mac and ip
type LocalMacIP struct {
	Mac string `json:"mac"`
	IP  string `json:"ip"`
}

func getRemoteIP() (string, error) {
	url := "https://srv1.futuredial.com/echosrv/echo/"
	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		// SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		fmt.Println("web request fail")
		return "", errors.New("web request fail")
	}
	// fmt.Println(string(resp.Body()))
	var dict map[string]map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &dict); err != nil {
		return "", errors.New("json format error")
	}
	if vv, ok := dict["request"]; ok {
		if v, ok := vv["remote_addr"]; ok {
			return v.(string), nil
		}
	}
	return "", errors.New("not find remote_addr")
}

// GetLocalPCInfo get local pc mac and ip
func GetLocalPCInfo() (string, string, error) {
	mac, ip, err := loadIPInfoFromFile()
	if err == nil {
		return mac, ip, nil
	}
	macip, err := Util.GetMapMacIP() // (map[string][]string, error)
	if err != nil {
		return "", "", nil
	}
	for mac, v := range macip {
		if len(v) > 0 {
			ip = v[0]
			savePCInfoToFile(mac, ip)
			return mac, ip, nil
		}

	}

	return "", "", errors.New("get mac ip failed")
}

func savePCInfoToFile(mac string, ip string) error {
	dat := map[string]string{
		mac: ip,
	}
	file, _ := json.MarshalIndent(dat, "", " ")
	return ioutil.WriteFile("localipinfo.json", file, 0644)
}

func loadIPInfoFromFile() (string, string, error) {
	jsonFile, err := os.Open("localipinfo.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		Log.Log.Error(err)
		return "", "", err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dat map[string]string
	if err := json.Unmarshal(byteValue, &dat); err != nil {
		// panic(err)
		Log.Log.Error(err)
		return "", "", err
	}
	var (
		mac string
		ip  string
	)
	for mac, ip = range dat {
		return mac, ip, nil
	}
	return "", "", errors.New("Not find mac ip from file")
}

// Download get download setting json
func Download() (map[string]interface{}, error) {
	configresult, err := dmc.GetSerialConfig()
	if err != nil {
		Log.Log.Error(err)
		return nil, err
	}
	pcname, _ := Util.GetPCName()
	mac, ip, _ := GetLocalPCInfo()
	remoteip, _ := getRemoteIP()
	var requestSetting RequestSetting
	requestSetting.Protocol = "3.0"
	companyid, _ := configresult.GetCompanyID()
	requestSetting.Client.Company = strconv.Itoa(companyid)
	requestSetting.Client.Solutionid = configresult.Solutionid
	requestSetting.Client.Productid = configresult.Productid
	requestSetting.Client.Site = configresult.Siteid
	requestSetting.Client.ID = configresult.ID
	requestSetting.Client.PCName = pcname
	requestSetting.Client.Macaddr = mac
	requestSetting.Client.Filecorrupted = "0"
	requestSetting.Client.Extrainfo.ID = configresult.ID
	requestSetting.Client.Extrainfo.IP = ip
	requestSetting.Client.Extrainfo.Product = configresult.PName
	requestSetting.Client.Extrainfo.RemoteAddr = remoteip

	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestSetting).
		Post(fmt.Sprintf("%s%s", configresult.Webserviceserver, "update/"))

	if err != nil {
		fmt.Println("web request fail")
		return nil, fmt.Errorf("Web request(get product settings)failed %v", err)
	}
	if resp.StatusCode() == http.StatusOK {
		var dat map[string]interface{}
		if err = json.Unmarshal(resp.Body(), &dat); err != nil {
			return nil, fmt.Errorf("data format error. %v", err)
		}
		return dat, nil
		// fmt.Println(string(resp.Body()))
	}
	return nil, fmt.Errorf("Web request(get product settings)failed %v", resp.Error())
}
