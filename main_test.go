package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
)

type Node struct {
	Text          string `xml:",chardata"`
	Index         string `xml:"index,attr"`
	AttrText      string `xml:"text,attr"`
	ResourceID    string `xml:"resource-id,attr"`
	Class         string `xml:"class,attr"`
	Package       string `xml:"package,attr"`
	ContentDesc   string `xml:"content-desc,attr"`
	Checkable     string `xml:"checkable,attr"`
	Checked       string `xml:"checked,attr"`
	Clickable     string `xml:"clickable,attr"`
	Enabled       string `xml:"enabled,attr"`
	Focusable     string `xml:"focusable,attr"`
	Focused       string `xml:"focused,attr"`
	Scrollable    string `xml:"scrollable,attr"`
	LongClickable string `xml:"long-clickable,attr"`
	Password      string `xml:"password,attr"`
	Selected      string `xml:"selected,attr"`
	Bounds        string `xml:"bounds,attr"`
	NAF           string `xml:"NAF,attr"`
	Node          []Node `xml:"node"`
}

type Hierarchy struct {
	XMLName  xml.Name `xml:"hierarchy"`
	Text     string   `xml:",chardata"`
	Rotation string   `xml:"rotation,attr"`
	Node     Node     `xml:"node"`
}

func TestXMLParser(t *testing.T) {
	xmls, err := os.Open("dump3.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Successfully Opened dump1.xml")
	// defer the closing of our jsonFile so that we can parse it later on
	defer xmls.Close()

	var v Hierarchy
	byteValue, _ := ioutil.ReadAll(xmls)
	err = xml.Unmarshal([]byte(byteValue), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(v)
}

func TestLogUDIDgetConfig(t *testing.T) {
	LogUDIDgetConfig("6c87ceb9-25a3-4e09-b81b-fb0a57b64d42")
}

func TestConfigParse(t *testing.T) {
	sss := `{
		"ok": 1, 
		"results": 
		[
			{
				"companyid": 83, 
				"staticfileserver": "http://cmcqa-dl.futuredial.com/", 
				"serverTime": "2020-01-22T23:17:51.679Z", 
				"pname": "CMC GreenT V5", 
				"installitunes": "False", 
				"adminconsoleserver": "http://cmcqa.futuredial.com/", 
				"siteid": "104", 
				"_id": "0d38e784-db8b-4ae2-8a96-e35e4f268240", 
				"solutionid": "1", 
				"webserviceserver": 
				"http://cmcqa.futuredial.com/ws/", 
				"productid": "1"
			}
		], 
		"id": 34524
		}`

	var dat cmc.ConfigInstall //map[string]interface{}
	if err := json.Unmarshal([]byte(sss), &dat); err != nil {
		panic(err)
	}
}

func TestParseLogResult(t *testing.T) {
	sss := `{
		"username": "qa", 
		"first_name": "", 
		"last_name": "", 
		"companygroup": null, 
		"last_pwd_reset": "2015-07-15T23:38:15+00:00", 
		"company": 9, 
		"is_active": true, 
		"site": 10, 
		"managedsites": [ ], 
		"email": "qa@futuredial.com", 
		"is_superuser": false, 
		"is_staff": false, 
		"last_login": "2020-07-09T17:17:17+00:00", 
		"available_products": [
			{
				"id": 9, 
				"name": "singlePC_GreenT"
			}, 
			{
				"id": 10, 
				"name": "singlePC_iRT"
			}, 
			{
				"id": 11, 
				"name": "TES"
			}, 
			{
				"id": 19, 
				"name": "GreenT"
			}, 
			{
				"id": 39, 
				"name": "Trust_iRT"
			}, 
			{
				"id": 51, 
				"name": "iTube"
			}, 
			{
				"id": 54, 
				"name": "allinone_golden_GreenT"
			}
		], 
		"groups": [
			4
		], 
		"crosssites": [ ], 
		"managedcompanys": [ ], 
		"id": 43, 
		"manageCredit": false, 
		"date_joined": "2015-07-15T23:38:15+00:00"
	}`
	var loginres cmc.LoginResult
	if err := json.Unmarshal([]byte(sss), &loginres); err != nil {
		panic(err)
	}
	ParseLogResult(loginres, 9, 9, 10, false)
}
