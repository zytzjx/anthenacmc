package cmcserverinfo

// ConfigResult
type ConfigResult struct {
	ID                 string `json:"_id"`
	Adminconsoleserver string `json:"adminconsoleserver"`
	Companyid          string `json:"companyid"`
	Installitunes      string `json:"installitunes"`
	PName              string `json:"pname"`
	ServerTime         string `json:"serverTime"`
	Staticfileserver   string `json:"staticfileserver"`
	Webserviceserver   string `json:"webserviceserver"`
	Productid          string `json:"productid"`
	Siteid             string `json:"siteid"`
	Solutionid         string `json:"solutionid"`
}

// ConfigInstall
type ConfigInstall struct {
	ID      int            `json:"id"`
	Ok      int            `json:"ok"`
	Results []ConfigResult `json:"results"`
}
