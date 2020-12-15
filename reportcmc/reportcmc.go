package reportcmc

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/juju/fslock"

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
	si, _ := config.GetSiteID()
	rbf.Site = strconv.Itoa(si)
	cid, _ := config.GetCompanyID()
	rbf.Company = strconv.Itoa(cid)
	pid, _ := config.GetProductID()
	rbf.Productid = strconv.Itoa(pid)
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

// Transcation send json to server
func Transcation(url string, info map[string]interface{}) (int, error) {
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
func ReportCMC(logfile string) (*ReportBaseFields, string, error) {
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
		sid, _ := configInstall.Results[0].GetSiteID()
		reportbase.Site = strconv.Itoa(sid)
		ii, _ := configInstall.Results[0].GetCompanyID()
		reportbase.Company = strconv.Itoa(ii)
		pid, _ := configInstall.Results[0].GetProductID()
		reportbase.Productid = strconv.Itoa(pid)
		reportbase.PortNumber = "1"
	}
	items, err := reportbase.MergeRedis()
	if err != nil {
		Log.Log.Error(err)
		return reportbase, staticurl, err
	}

	_, err = Transcation(configInstall.Results[0].Webserviceserver, items)
	if err != nil {
		Log.Log.Error(err)
		saveDatatoFile(items, logfile)
		return reportbase, staticurl, err
	}
	return reportbase, staticurl, nil
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func saveDatatoFile(items map[string]interface{}, logfile string) error {
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
	err = ioutil.WriteFile("transcationpool/"+uuid+".json", file, 0644)
	if err != nil {
		return err
	}
	if _, err := os.Stat(logfile); err == nil {
		Log.Log.Info("copy log to backup")
		copyFileContents(logfile, "transcationpool/"+uuid+".zip")
	}
	return nil
}

// PostLogFile to server
func PostLogFile(url, uuid, productid string, filePath string) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Log.Log.Infof("%s file not exist", filePath)
		return nil
	}

	//打开文件句柄操作
	file, err := os.Open(filePath)
	if err != nil {
		Log.Log.Error("error opening file")
		return err
	}
	defer file.Close()

	//创建一个模拟的form中的一个选项,这个form项现在是空的
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作, 设置文件的上传参数叫uploadfile, 文件名是filename,
	//相当于现在还没选择文件, form项里选择文件的选项
	fileWriter, err := bodyWriter.CreateFormFile("fileobj", uuid+".zip")
	if err != nil {
		Log.Log.Error("error writing to buffer")
		return err
	}

	//iocopy 这里相当于选择了文件,将文件放到form中
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	//获取上传文件的类型,multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()

	//这里就是上传的其他参数设置,可以使用 bodyWriter.WriteField(key, val) 方法
	//也可以自己在重新使用  multipart.NewWriter 重新建立一项,这个再server 会有例子
	params := map[string]string{
		"uuid":      uuid,
		"productid": productid,
	}
	//这种设置值得仿佛 和下面再从新创建一个的一样
	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}
	//这个很关键,必须这样写关闭,不能使用defer关闭,不然会导致错误
	bodyWriter.Close()

	//发送post请求到服务端
	url = fmt.Sprintf("%s%s", url, "uploadlog/")
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respbody))
	// ioutil.WriteFile("fileresult.txt", respbody, 0644)

	if resp.StatusCode == http.StatusOK {
		rsjs := make(map[string]interface{})
		if json.Unmarshal(respbody, &rsjs) != nil {
			if errstr, ok := rsjs["error"]; ok {
				return fmt.Errorf("http error, %s", errstr)
			}
		}
		return nil
	}
	return fmt.Errorf("http error, %s", resp.Status)

}

// SendLocalFiletoCMC send file to cmc
func SendLocalFiletoCMC(serviceserver string, staticserver string) {
	// var configInstall cmc.ConfigInstall //map[string]interface{}
	// if err := configInstall.LoadFile("serialconfig.json"); err != nil {
	// 	Log.Log.Error(err)
	// 	return
	// }
	lock := fslock.New(".uploadlist.lock")
	lockErr := lock.TryLock()
	if lockErr != nil {
		return
	}
	files, err := ioutil.ReadDir("transcationpool")
	// release the lock
	lock.Unlock()

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
			uuid := fmt.Sprintf("%v", items["uuid"])
			productid := fmt.Sprintf("%v", items["productid"])
			logfile := filepath.Join("transcationpool", uuid+".zip")

			_, err = Transcation(serviceserver, items)
			if err == nil {
				if PostLogFile(staticserver, uuid, productid, logfile) != nil {
					continue
				}
				os.Remove(filepath.Join("transcationpool", file.Name()))
				os.Remove(logfile)
			}
		}
	}
}
