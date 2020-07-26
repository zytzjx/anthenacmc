package ftpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	Log "github.com/zytzjx/anthenacmc/loggersys"
)

var (
	// FTPSERVER ftp server
	FTPSERVER string = "ftp2.futuredial.com"
	// FTPPORT ftp port, default 21
	FTPPORT int = 21
	// USERNAME username for ftp
	USERNAME string = "fd_eng"
	// PASSWORD ftp password
	PASSWORD string = "FDeng"
	// ROOTDIR root dir
	ROOTDIR string = "/ModusLink/CMC_Report"
)

// FTPINFO log ftp
type FTPINFO struct {
	FTPServer string `json:"server"`
	FTPPort   int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	RootDir   string `json:"rootdir"`
}

func newFTPInfo() *FTPINFO {
	p := FTPINFO{}
	p.FTPServer = FTPSERVER
	p.FTPPort = FTPPORT
	p.User = USERNAME
	p.Password = PASSWORD
	p.RootDir = ROOTDIR
	return &p
}

// LoadConfig load config
func (fi *FTPINFO) LoadConfig(conf string) error {
	// Open our jsonFile
	jsonFile, err := os.Open(conf)
	// if we os.Open returns an error then handle it
	if err != nil {
		Log.Log.Error(err)
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, fi); err != nil {
		// panic(err)
		Log.Log.Error(err)
		return err
	}
	return err
}

// PrintInfo print to log
func (fi *FTPINFO) PrintInfo() {
	fmt.Println(fi)
}

// CreateFolder create ftp folder
func (fi *FTPINFO) CreateFolder(folder string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", fi.FTPServer, fi.FTPPort), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	if err = c.Login(fi.User, fi.Password); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(fi.RootDir, folder)
	err = c.MakeDir(path)
	return err
}

// RemoveFolder remove ftp folder
func (fi *FTPINFO) RemoveFolder(folder string, bRecur bool) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", fi.FTPServer, fi.FTPPort), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	if err = c.Login(fi.User, fi.Password); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(fi.RootDir, folder)
	if bRecur {
		err = c.RemoveDirRecur(path)
	} else {
		err = c.RemoveDir(path)
	}

	return err
}

// Upload upload local file to ftp file
func (fi *FTPINFO) Upload(localname, filename string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", fi.FTPServer, fi.FTPPort), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	if err = c.Login(fi.User, fi.Password); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(fi.RootDir, filename)
	var file *os.File
	if file, err = os.Open(localname); err != nil {
		Log.Log.Fatal(err)
		return err
	}

	err = c.Stor(path, file)

	return err
}

// Download from ftp
func (fi *FTPINFO) Download(localname, filename string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", fi.FTPServer, fi.FTPPort), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	if err = c.Login(fi.User, fi.Password); err != nil {
		Log.Log.Fatal(err)
		return err
	}

	path := filepath.Join(fi.RootDir, filename)
	var file *os.File
	if file, err = os.Create(localname); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer file.Close()

	r, err := c.Retr(path)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	file.Write(buf)
	return err
}

// DeleteFile delete file
func (fi *FTPINFO) DeleteFile(filename string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", fi.FTPServer, fi.FTPPort), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	if err = c.Login(fi.User, fi.Password); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(fi.RootDir, filename)
	return c.Delete(path)
}
