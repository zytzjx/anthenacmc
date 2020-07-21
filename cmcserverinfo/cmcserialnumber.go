package cmcserverinfo

import (
	"errors"
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
