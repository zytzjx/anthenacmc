package reportcmc

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
	"github.com/zytzjx/anthenacmc/datacentre"
	Log "github.com/zytzjx/anthenacmc/loggersys"
	"github.com/zytzjx/anthenacmc/utils"
)

// ReportBaseFields MUST INCLUDE Field
type ReportBaseFields struct {
	ID          string `json:"_id"`
	UUID        string `json:"uuid"`
	Site        string `json:"site"`
	Company     string `json:"company"`
	Productid   string `json:"productid"`
	TimeCreated string `json:"timeCreated"`
	EsnNumber   string `json:"esnNumber"`
	PortNumber  string `json:"portNumber"`
	Operator    string `json:"operator"`
	SourceModel string `json:"sourceModel"`
	SourceMake  string `json:"sourceMake"`
	ErrorCode   string `json:"errorCode"`
}

// NewReportBaseFields create report MUST field
func NewReportBaseFields() *ReportBaseFields {
	rbf := ReportBaseFields{}
	id, err := newUUID()
	if err != nil {
		Log.Log.Error("new uuid failed")
		return nil
	}
	rbf.ID = id
	rbf.UUID = id
	rbf.TimeCreated = time.Now().Format("2006-01-02T15:04:05.000000000")

	config, err := datacentre.GetSerialConfig()
	if err != nil {
		Log.Log.Error("GetSerialConfig failed")
		return &rbf
	}
	rbf.Site = config.Siteid
	rbf.Company = config.Companyid.(string)
	rbf.Productid = config.Productid
	rbf.PortNumber = "1"
	rbf.ErrorCode = "1"
	rbf.EsnNumber = "000000000000000"
	rbf.SourceMake = "Futuredial"
	rbf.Operator = "00000"
	rbf.SourceModel = "PST_GRADING_FD"

	return &rbf
}

// SetDevice set Device
func (rbf *ReportBaseFields) SetDevice(esn, model, make, errorcode string) {
	rbf.EsnNumber = esn
	rbf.SourceMake = make
	rbf.SourceModel = model
	rbf.ErrorCode = errorcode
}

// MergeRedis to cmc server
func (rbf *ReportBaseFields) MergeRedis() (map[string]interface{}, error) {
	data := utils.StructToMap(rbf)
	info, err := datacentre.GetTransaction()
	if err != nil {
		Log.Log.Error(err)
		return nil, err
	}
	for k, v := range info {
		data[k] = v
	}
	return data, nil
}

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

func transcation(url string, info map[string]interface{}) (int, error) {
	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(info).
		Post(fmt.Sprintf("%s%s", url, "insert/"))

	if err != nil {
		fmt.Println("web request fail")
		Log.Log.Errorf("web request fail: %s", err)
		return 0, err
	}
	Log.Log.Info(string(resp.Body()))

	var status map[string]interface{}
	// fmt.Println(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), &status); err != nil {
		// fmt.Println("return format error.")
		Log.Log.Error("return format error.")
		return 0, err
	}

	if val, ok := status["status"]; ok {
		var stat int
		switch val.(type) {
		case string:
			stat, _ = strconv.Atoi(val.(string))
		case float64:
			stat = int(val.(float64))
		case int, int32, int64:
			stat = int(val.(int))
		}
		switch stat {
		case 1:
			return 1, nil
		case 2:
			// fmt.Println("success")
			Log.Log.Info("Success")
			return 2, nil
		case 3:
			return 3, fmt.Errorf("%v", status["error"])
		case 4:
			Log.Log.Warn("same uuid")
			// fmt.Println("same uuid")
			return 4, nil
		default:
			// fmt.Println("failed.")
			Log.Log.Error("failed.")
		}

	}
	return 0, errors.New("I do not know the protocol return")
}

// ReportCMC to server
func ReportCMC() (*ReportBaseFields, string, error) {
	var staticurl string
	var configInstall cmc.ConfigInstall //map[string]interface{}
	if err := configInstall.LoadFile("serialconfig.json"); err != nil {
		Log.Log.Error(err)
		return nil, staticurl, err
	}
	staticurl = configInstall.Results[0].Staticfileserver
	reportbase := NewReportBaseFields()
	if reportbase == nil {
		err := errors.New("data base create failed")
		Log.Log.Error(err)
		return reportbase, staticurl, err
	}

	reportbase.Operator, _ = datacentre.GetString("login.operator")

	if reportbase.Company == "" || reportbase.PortNumber == "" {
		reportbase.Site = configInstall.Results[0].Siteid
		ii, _ := configInstall.Results[0].GetCompanyID()
		reportbase.Company = strconv.Itoa(ii)
		reportbase.Productid = configInstall.Results[0].Productid
		reportbase.PortNumber = "1"
	}
	items, err := reportbase.MergeRedis()
	if err != nil {
		Log.Log.Error(err)
		return reportbase, staticurl, err
	}

	_, err = transcation(configInstall.Results[0].Webserviceserver, items)
	if err != nil {
		Log.Log.Error(err)
		saveDatatoFile(items)
		return reportbase, staticurl, err
	}
	return reportbase, staticurl, nil
}

func saveDatatoFile(items map[string]interface{}) error {
	uuid := fmt.Sprintf("%v", items["uuid"])
	if _, err := os.Stat("transcationpool"); os.IsNotExist(err) {
		// /var/log/anthena does not exist
		if err = os.Mkdir("transcationpool", 0775); err != nil {
			fmt.Println(err)
		}
	}
	file, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(uuid+".json", file, 0644)
	return err
}

// SendLocalFiletoCMC send file to cmc
func SendLocalFiletoCMC() {
	var configInstall cmc.ConfigInstall //map[string]interface{}
	if err := configInstall.LoadFile("serialconfig.json"); err != nil {
		Log.Log.Error(err)
		return
	}
	files, err := ioutil.ReadDir("transcationpool")
	if err != nil {
		Log.Log.Error(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".json" {
			jsonFile, err := os.Open(filepath.Join("transcationpool", file.Name()))
			// if we os.Open returns an error then handle it
			if err != nil {
				Log.Log.Error(err)
				continue
			}
			// defer the closing of our jsonFile so that we can parse it later on
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)
			var items map[string]interface{}
			if err := json.Unmarshal(byteValue, &items); err != nil {
				Log.Log.Error(err)
				continue
			}
			_, err = transcation(configInstall.Results[0].Webserviceserver, items)
			if err == nil {
				os.Remove(filepath.Join("transcationpool", file.Name()))
			}
		}
	}
}
