package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"

	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
	dmc "github.com/zytzjx/anthenacmc/datacentre"
	Log "github.com/zytzjx/anthenacmc/loggersys"
	_ "github.com/zytzjx/anthenacmc/loggersys"
)

// User login usename and password
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LogUDIDgetConfig  dd
func LogUDIDgetConfig(uuid string) int {
	// https://ps.futuredial.com/profiles/clients/_find?criteria={"_id":"ed2e7151-441d-4f42-9916-7794a55abb0e"}
	client := resty.New()
	// 196003c3-5c93-4c98-94d9-96e945122a5d
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"criteria": fmt.Sprintf("{\"_id\":\"%s\"}", uuid),
			// "criteria": "{\"_id\":\"0d38e784-db8b-4ae2-8a96-e35e4f268240\"}",
		}).
		SetHeader("Accept", "application/json").
		Get("https://ps.futuredial.com/profiles/clients/_find")
	if err != nil {
		fmt.Println("web request failed.")
		Log.Log.Error("web request failed.")
		return 1
	}

	var dat cmc.ConfigInstall //map[string]interface{}
	Log.Log.Info(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), &dat); err != nil {
		// panic(err)
		Log.Log.Error(err)
		return 2
	}
	dmc.SaveSerialConfigRedis(dat)
	file, _ := json.MarshalIndent(dat, "", " ")
	_ = ioutil.WriteFile("serialconfig.json", file, 0644)

	return 0
	// loginCMC(dat, User{Username: "qa", Password: "qa"})
	// for k, v := range dat {
	// 	switch vv := v.(type) {
	// 	case string:
	// 		fmt.Println(k, "is string", vv)
	// 	case int64:
	// 		fmt.Println(k, "is int", vv)
	// 	case bool:
	// 		fmt.Println(k, "is Bool", vv)
	// 	case float64:
	// 		fmt.Println(k, "is float64", vv)
	// 	case []interface{}:
	// 		fmt.Println(k, "is an array:")
	// 		for _, u := range vv {
	// 			mm := u.(map[string]interface{})
	// 			if sbinterface, bok := dat["installitunes"]; bok {
	// 				fmt.Println(sbinterface)
	// 				if sb, bbok := sbinterface.(bool); bbok {
	// 					if sb {
	// 						fmt.Println("installitunes", sb)
	// 					} else {
	// 						fmt.Println("installitunes", sb)
	// 					}
	// 				}
	// 			}
	// 			for kk, vk := range mm {
	// 				fmt.Println(kk, "===", vk)
	// 				if str, ok := vk.(string); ok {
	// 					// if !strings.EqualFold(kk, "_id") {
	// 					// 	WritePrivateProfileString("config", kk, str, filename)
	// 					// }
	// 					fmt.Println(str)
	// 				} else {
	// 					fmt.Println("format error!")
	// 				}
	// 			}
	// 		}
	// 	default:
	// 		fmt.Println(k, "is of a type I don't know how to handle")
	// 	}
	// }
}

// loginCMC
func loginCMC(config cmc.ConfigInstall, usr User) (*cmc.LoginResult, error) {
	if config.Ok != 1 {
		return nil, errors.New("serial config is not sucessful")
	}
	url := config.Results[0].Adminconsoleserver

	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(usr).
		Post(fmt.Sprintf("%s%s", url, "api/auth/"))

	if err != nil {
		fmt.Println("web request fail")
		return nil, errors.New("web request fail")
	}
	var dict map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &dict); err != nil {
		return nil, errors.New("json format error")
	}
	if err, ok := dict["error"]; ok {
		fmt.Println(err)
		//return 1 //error
		return nil, errors.New(err.(string))
	}

	fmt.Println(string(resp.Body()))
	var loginres cmc.LoginResult
	if err := json.Unmarshal(resp.Body(), &loginres); err != nil {
		return nil, errors.New("json format parse to Login Result error")
	}
	// for _, u := range loginres.AvailableProducts {
	// 	fmt.Println(u)
	// }
	return &loginres, nil

}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// ParseLogResult parser log in
// string company_id = companyId();
// string product_id = productId();
// string site_id = siteId();
// string strcheckid = checksiteId();
func ParseLogResult(loginres cmc.LoginResult, companyID int, productID int, siteID int, allowCheckSiteID bool) int {
	var operatorvalue bool
	var bProduct bool
	ret := 0

	if _, ok := Find(loginres.Groups, 4); ok {
		operatorvalue = true
	}

	for _, u := range loginres.AvailableProducts {
		if u.ID == productID {
			bProduct = true
			break
		}
	}
	if allowCheckSiteID {
		if loginres.IsActive && siteID == loginres.Site && companyID == loginres.Company && operatorvalue && bProduct {
			//save xml
			file, _ := json.MarshalIndent(loginres, "", " ")
			_ = ioutil.WriteFile("hydralogin.json", file, 0644)
		} else if loginres.IsActive && companyID == loginres.Company && operatorvalue && bProduct {
			return ret
		} else {
			if !loginres.IsActive {
				ret = 5
			} else if companyID != loginres.Company {
				ret = 3
			} else if siteID != loginres.Site {
				ret = 9
			} else if !operatorvalue {
				ret = 2
			} else if !bProduct {
				ret = 7
			}
		}
	} else {
		if loginres.IsActive && companyID == loginres.Company && operatorvalue && bProduct {
			//xml.Save(Path.Combine(Program.output, "hydralogin.xml"));
			file, _ := json.MarshalIndent(loginres, "", " ")
			_ = ioutil.WriteFile("hydralogin.json", file, 0644)
		} else if loginres.IsActive && companyID == loginres.Company && !operatorvalue && bProduct {
			ret = 8
		} else {
			if !loginres.IsActive {
				ret = 5
			} else if companyID != loginres.Company {
				ret = 3
			} else if !operatorvalue {
				ret = 2
			} else if !bProduct {
				ret = 7
			}
		}
	}
	return ret
}

func main() {
	bLogin := flag.Bool("login", false, "-login login or get project config by serialnamber")
	uuid := flag.String("uuid", "", "serialnumber of the project ")
	username := flag.String("username", "", "login user name")
	password := flag.String("password", "", "login password")

	flag.Parse()

	var ret int
	if *uuid != "" {
		ret = LogUDIDgetConfig(*uuid)
		if ret != 0 {
			os.Exit(ret)
		}
	}

	if *bLogin {
		// Open our jsonFile
		jsonFile, err := os.Open("serialconfig.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			Log.Log.Error(err)
			os.Exit(4)
		}
		fmt.Println("Successfully Opened serialconfig.json")
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var dat cmc.ConfigInstall //map[string]interface{}
		if err := json.Unmarshal(byteValue, &dat); err != nil {
			// panic(err)
			Log.Log.Error(err)
			os.Exit(3)
		}

		loginres, err := loginCMC(dat, User{Username: *username, Password: *password})
		if err != nil {
			Log.Log.Error(err)
			os.Exit(5)
		}
		Companyid, _ := dat.Results[0].GetCompanyID()
		Productid, _ := strconv.Atoi(dat.Results[0].Productid)
		Siteid, _ := strconv.Atoi(dat.Results[0].Siteid)
		ret = ParseLogResult(*loginres, Companyid, Productid, Siteid, false)
	}
	os.Exit(ret)

}
