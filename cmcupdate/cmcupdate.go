package cmcupdate

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/go-resty/resty/v2"
	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
	Log "github.com/zytzjx/anthenacmc/loggersys"
	"github.com/zytzjx/anthenacmc/utils"
)

// ClientInfo cmc UUID  response
type ClientInfo struct {
	Company    int `json:"company"`
	Solutionid int `json:"solutionid"`
	Productid  int `json:"productid"`
}

// ModuleFiles phonedll and phonetips
type ModuleFiles struct {
	Filelist   []map[string]interface{} `json:"filelist"`
	Deletelist []map[string]interface{} `json:"deletelist"`
}

func newModuleFiles() *ModuleFiles {
	mf := ModuleFiles{}
	mf.Deletelist = make([]map[string]interface{}, 0)
	mf.Filelist = make([]map[string]interface{}, 0)
	return &mf
}

// FrameworkFiles framework info
type FrameworkFiles struct {
	Version    string                   `json:"version"`
	Filelist   []map[string]interface{} `json:"filelist"`
	Deletelist []map[string]interface{} `json:"deletelist"`
}

// SyncStatus request json
type SyncStatus struct {
	Client ClientInfo `json:"client"`
	Sync   struct {
		Status struct {
			Framework FrameworkFiles `json:"framework"`
			Phonedll  ModuleFiles    `json:"phonedll"`
			Phonetips ModuleFiles    `json:"phonetips"`
			// settings map[string]interface{} `json:"settings"`
		} `json:"status"`
	} `json:"sync"`
	Protocol string `json:"protocol"`
}

// StatusResponse server response for download
type StatusResponse struct {
	Framework FrameworkFiles         `json:"framework"`
	Phonedll  ModuleFiles            `json:"phonedll"`
	Phonetips ModuleFiles            `json:"phonetips"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
}

// newSyncStatus for first request
func newSyncStatus(cliinfo ClientInfo) *SyncStatus {
	sync := SyncStatus{}
	sync.Protocol = "2.0"
	sync.Client = cliinfo
	sync.Sync.Status.Framework.Filelist = make([]map[string]interface{}, 0)
	sync.Sync.Status.Phonedll = *newModuleFiles()
	sync.Sync.Status.Phonetips = *newModuleFiles()
	return &sync
}

// ModuleFileItem is ModuleFiles {List } items
type ModuleFileItem struct {
	Checksum    string      `json:"checksum"`
	Disabled    bool        `json:"disabled,omitempty"`
	Readableid  string      `json:"readableid,omitempty"`
	Size        interface{} `json:"size"`
	DownloadURL string      `json:"url,omitempty"`
}

// GetFileSize Size field sametime is string , sometime is string
func (mfi ModuleFileItem) GetFileSize() (int, error) {
	var size int
	var err error
	switch mfi.Size.(type) {
	case int, int64, int16, int32:
		size = mfi.Size.(int)
		return size, nil
	case string:
		size, err = strconv.Atoi(mfi.Size.(string))
		return size, err
	default:
		err = errors.New("no support format")
		return 0, err
	}
}

// saveStatusFile save for download
func saveStatusFile(jsondata []byte) error {
	return ioutil.WriteFile("res_syncstatus.json", jsondata, 0644)
}

// sendRequest send request to cmc server
func sendRequest(url string, syncstauts SyncStatus) (StatusResponse, error) {
	var download StatusResponse
	// ss, err := json.Marshal(syncstauts)
	// if err != nil {
	// 	return download, errors.New("request json format error")
	// }
	// fmt.Println(string(ss))
	// Create a Resty Client
	client := resty.New()
	// POST Struct, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(syncstauts).
		Post(url)

	if err != nil {
		fmt.Println("web request fail")
		return download, errors.New("web request fail")
	}
	// Log.Log.Info(string(resp.Body()))
	var dict map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &dict); err != nil {
		return download, errors.New("json format error")
	}
	if err, ok := dict["error"]; ok {
		fmt.Println(err)
		//return 1 //error
		return download, errors.New(err.(string))
	}

	// var download StatusResponse
	if err := json.Unmarshal(resp.Body(), &download); err != nil {
		return download, errors.New("json format is not protocol error")
	}

	// fileitem := download.Phonedll.Filelist
	// fmt.Println(fileitem)
	// file, err := json.MarshalIndent(fileitem, "", " ")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println((string)(file))
	//  var mlst []ModuleFileItem
	// json.Unmarshal(file, &mlst)
	// fmt.Println(mlst[0].updateURL)
	saveStatusFile(resp.Body())
	return download, nil
}

// loadcompanysetting from file
func loadcompanysetting() (ClientInfo, string, error) {
	// Open our jsonFile
	var cliinfo ClientInfo

	jsonFile, err := os.Open("serialconfig.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		Log.Log.Error(err)
		return cliinfo, "", fmt.Errorf("serial number run first. %s", err)
	}
	fmt.Println("Successfully Opened serialconfig.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dat cmc.ConfigInstall //map[string]interface{}
	if err := json.Unmarshal(byteValue, &dat); err != nil {
		// panic(err)
		return cliinfo, "", err
	}
	cliinfo.Company, _ = dat.Results[0].GetCompanyID()
	cliinfo.Productid, _ = strconv.Atoi(dat.Results[0].Productid)
	cliinfo.Solutionid, _ = strconv.Atoi(dat.Results[0].Solutionid)
	updateurl := dat.Results[0].Webserviceserver

	return cliinfo, updateurl, nil

}

var syncstatus SyncStatus

// UpdateCMC update data from CMC
func UpdateCMC() (StatusResponse, error) {
	cliinfo, updateurl, err := loadcompanysetting()
	if err != nil {
		Log.Log.Error(err)
		return StatusResponse{}, err
	}
	updateurl = updateurl + "update/"

	// clientstatus
	jsonFile, err := os.Open("clientstatus.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		Log.Log.Error(err)
		syncstatus = *newSyncStatus(cliinfo)
		return sendRequest(updateurl, syncstatus)
	}
	fmt.Println("Successfully Opened serialconfig.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &syncstatus); err != nil {
		return StatusResponse{}, err
	}

	return sendRequest(updateurl, syncstatus)
}

func httpdownload(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func md5file(localpath, checksum string) (bool, error) {
	f, err := os.Open(localpath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	md5file := fmt.Sprintf("%X", h.Sum(nil))
	if md5file != checksum {
		return false, errors.New("download file checksum failed")
	}
	return true, nil
}

func downloadFile(mfi ModuleFileItem) error {
	file := path.Base(mfi.DownloadURL)
	localpath := path.Join("update", file)

	if utils.FileExists(localpath) {
		if ok, _ := md5file(localpath, mfi.Checksum); ok {
			return nil
		}
	}

	if err := httpdownload(mfi.DownloadURL, localpath); err != nil {
		return err
	}

	f, err := os.Open(localpath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	md5file := fmt.Sprintf("%X", h.Sum(nil))
	if md5file != mfi.Checksum {
		return errors.New("download file checksum failed")
	}
	return nil
}

// changeToModuleItems
func changeToModuleItems(filelist []map[string]interface{}) ([]ModuleFileItem, error) {
	file, err := json.MarshalIndent(filelist, "", " ")
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(file))
	var mlst []ModuleFileItem
	if err = json.Unmarshal(file, &mlst); err != nil {
		return nil, err
	}
	return mlst, nil
}

// DownloadCMC dowload from CMC server
func DownloadCMC() ([]ModuleFileItem, error) {
	var faildowndlist []ModuleFileItem
	if _, err := os.Stat("update"); os.IsNotExist(err) {
		// /var/log/anthena does not exist
		if err = os.Mkdir("update", 0775); err != nil {
			Log.Log.Error(err)
			return faildowndlist, err
		}
	}
	sr, err := UpdateCMC()
	if err != nil {
		Log.Log.Error(err)
		return faildowndlist, err
	}
	count := 0
	frameworks, err1 := changeToModuleItems(sr.Framework.Filelist)
	if err1 == nil {
		count += len(frameworks)
	}
	phonedll, err2 := changeToModuleItems(sr.Phonedll.Filelist)
	if err2 == nil {
		count += len(phonedll)
	}
	phonetip, err3 := changeToModuleItems(sr.Phonetips.Filelist)
	if err3 == nil {
		count += len(phonetip)
	}
	var wg sync.WaitGroup
	wg.Add(count)
	queue := make(chan ModuleFileItem, 1)
	if err1 == nil {
		for _, it := range frameworks {
			go func(mfi ModuleFileItem, wg *sync.WaitGroup) error {
				err := downloadFile(mfi)
				if err != nil {
					queue <- mfi
				}
				wg.Done()
				return err
			}(it, &wg)
		}
	}
	if err2 == nil {
		for _, it := range phonedll {
			go func(mfi ModuleFileItem, wg *sync.WaitGroup) error {
				err := downloadFile(mfi)
				if err != nil {
					queue <- mfi
				}
				wg.Done()
				return err
			}(it, &wg)
		}
	}
	if err3 == nil {
		for _, it := range phonetip {
			go func(mfi ModuleFileItem, wg *sync.WaitGroup) error {
				err := downloadFile(mfi)
				if err != nil {
					queue <- mfi
				}
				wg.Done()
				return err
			}(it, &wg)
		}
	}

	go func() {
		for t := range queue {
			faildowndlist = append(faildowndlist, t)
		}
	}()

	wg.Wait()
	fmt.Println(faildowndlist)
	return faildowndlist, nil
}

// RetryDownload download fail
func RetryDownload(items []ModuleFileItem) ([]ModuleFileItem, error) {
	var wg sync.WaitGroup
	count := len(items)
	wg.Add(count)
	var faildowndlist []ModuleFileItem
	queue := make(chan ModuleFileItem, 1)
	for _, it := range items {
		go func(mfi ModuleFileItem, wg *sync.WaitGroup) error {
			err := downloadFile(mfi)
			if err != nil {
				queue <- mfi
			}
			wg.Done()
			return err
		}(it, &wg)
	}

	go func() {
		for t := range queue {
			faildowndlist = append(faildowndlist, t)
		}
	}()

	wg.Wait()
	if len(faildowndlist) > 0 {
		return faildowndlist, errors.New("download failed")
	}
	return faildowndlist, nil

}
