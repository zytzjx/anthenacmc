package reportcmc

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-resty/resty/v2"

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
		return nil
	}
	rbf.ID = id
	rbf.UUID = id
	rbf.TimeCreated = time.Now().Format("2006-01-02T15:04:05.000000000")

	config, err := datacentre.GetSerialConfig()
	if err != nil {
		return nil
	}
	rbf.Site = config.Siteid
	rbf.Company = config.Companyid.(string)
	rbf.Productid = config.Productid
	rbf.PortNumber = "1"

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

func transcation(url string, info map[string]interface{}) {
	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(info).
		Post(fmt.Sprintf("%s%s", url, "insert/"))

	if err != nil {
		fmt.Println("web request fail")
		Log.Log.Error("web request fail")
		return
	}
	Log.Log.Info(string(resp.Body()))

	var status map[string]interface{}
	// fmt.Println(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), &status); err != nil {
		// fmt.Println("return format error.")
		Log.Log.Error("return format error.")
	}

	if val, ok := status["status"]; ok {
		switch val {
		case "1", "2":
			// fmt.Println("success")
			Log.Log.Info("Success")
		case "4":
			Log.Log.Warn("same uuid")
			// fmt.Println("same uuid")
		default:
			// fmt.Println("failed.")
			Log.Log.Error("failed.")
		}
	}

}
