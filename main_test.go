package main

import (
	"encoding/json"
	"testing"

	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
)

func TestLogUDIDgetConfig(t *testing.T) {
	LogUDIDgetConfig("6c87ceb9-25a3-4e09-b81b-fb0a57b64d42")
}

func TestConfigParse(t *testing.T) {
	sss := `{
		"ok": 1, 
		"results": 
		[
			{
				"companyid": "83", 
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
