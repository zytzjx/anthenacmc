package cmcserverinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// ConfigResult
type ConfigResult struct {
	ID                 string      `json:"_id"`
	Adminconsoleserver string      `json:"adminconsoleserver"`
	Companyid          interface{} `json:"companyid"`
	Installitunes      string      `json:"installitunes"`
	PName              string      `json:"pname"`
	ServerTime         string      `json:"serverTime"`
	Staticfileserver   string      `json:"staticfileserver"`
	Webserviceserver   string      `json:"webserviceserver"`
	Productid          string      `json:"productid"`
	Siteid             string      `json:"siteid"`
	Solutionid         string      `json:"solutionid"`
}

// ConfigInstall uuid return
type ConfigInstall struct {
	ID      int            `json:"id"`
	Ok      int            `json:"ok"`
	Results []ConfigResult `json:"results"`
}

// GetCompanyID company id , because server return sometime is int , sometimes is string
func (cr ConfigResult) GetCompanyID() (int, error) {
	var Companyid int
	var err error
	switch cr.Companyid.(type) {
	case int, int64, int16, int32:
		Companyid = cr.Companyid.(int)
		return Companyid, nil
	case string:
		Companyid, err = strconv.Atoi(cr.Companyid.(string))
		return Companyid, err
	default:
		err = errors.New("no support format")
		return 0, err
	}
}

// LoadFile load serialconfig.json file
func (ci *ConfigInstall) LoadFile(filename string) error {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	fmt.Println("Successfully Opened serialconfig.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	// var dat cmc.ConfigInstall //map[string]interface{}
	if err := json.Unmarshal(byteValue, ci); err != nil {
		// panic(err)
		return err
	}
	return nil
}
