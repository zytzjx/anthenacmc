package cmcserverinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// ConfigResult from server
type ConfigResult struct {
	ID                 string      `json:"_id"`
	Adminconsoleserver string      `json:"adminconsoleserver"`
	Companyid          interface{} `json:"companyid"`
	Installitunes      string      `json:"installitunes"`
	PName              string      `json:"pname"`
	ServerTime         string      `json:"serverTime"`
	Staticfileserver   string      `json:"staticfileserver"`
	Webserviceserver   string      `json:"webserviceserver"`
	Productid          interface{} `json:"productid"`
	Siteid             interface{} `json:"siteid"`
	Solutionid         interface{} `json:"solutionid"`
}

// ConfigInstall uuid return
type ConfigInstall struct {
	ID      int            `json:"id"`
	Ok      int            `json:"ok"`
	Results []ConfigResult `json:"results"`
}

// GetProductID , because server return sometime is int , sometimes is string
func (cr ConfigResult) GetProductID() (int, error) {
	var productid int
	var err error
	switch cr.Productid.(type) {
	case int, int64, int16, int32:
		productid = cr.Productid.(int)
		return productid, nil
	case string:
		productid, err = strconv.Atoi(cr.Productid.(string))
		return productid, err
	default:
		err = errors.New("no support format")
		return 0, err
	}
}

// GetSiteID , because server return sometime is int , sometimes is string
func (cr ConfigResult) GetSiteID() (int, error) {
	var siteid int
	var err error
	switch cr.Siteid.(type) {
	case int, int64, int16, int32:
		siteid = cr.Siteid.(int)
		return siteid, nil
	case string:
		siteid, err = strconv.Atoi(cr.Siteid.(string))
		return siteid, err
	default:
		err = errors.New("no support format")
		return 0, err
	}
}

// GetSolutionID company id , because server return sometime is int , sometimes is string
func (cr ConfigResult) GetSolutionID() (int, error) {
	var Solutionid int
	var err error
	switch cr.Solutionid.(type) {
	case int, int64, int16, int32:
		Solutionid = cr.Solutionid.(int)
		return Solutionid, nil
	case string:
		Solutionid, err = strconv.Atoi(cr.Solutionid.(string))
		return Solutionid, err
	default:
		err = errors.New("no support format")
		return 0, err
	}
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
