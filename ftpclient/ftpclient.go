package ftpclient

import (
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

// CreateFolder create ftp folder
func CreateFolder(folder string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", FTPSERVER, FTPPORT), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	err = c.Login(USERNAME, PASSWORD)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(ROOTDIR, folder)
	err = c.MakeDir(path)
	return err
}

// RemoveFolder remove ftp folder
func RemoveFolder(folder string, bRecur bool) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", FTPSERVER, FTPPORT), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	err = c.Login(USERNAME, PASSWORD)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(ROOTDIR, folder)
	if bRecur {
		err = c.RemoveDirRecur(path)
	} else {
		err = c.RemoveDir(path)
	}

	return err
}

// Upload upload local file to ftp file
func Upload(localname, filename string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", FTPSERVER, FTPPORT), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	err = c.Login(USERNAME, PASSWORD)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(ROOTDIR, filename)
	var file *os.File
	if file, err = os.Open(localname); err != nil {
		Log.Log.Fatal(err)
		return err
	}

	err = c.Stor(path, file)

	return err
}

func Download(localname, filename string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", FTPSERVER, FTPPORT), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer c.Quit()

	err = c.Login(USERNAME, PASSWORD)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	path := filepath.Join(ROOTDIR, filename)
	var file *os.File
	if file, err = os.Create(localname); err != nil {
		Log.Log.Fatal(err)
		return err
	}
	resp, err = c.Retr(filename)
	if err != nil {
		Log.Log.Fatal(err)
		return err
	}
	defer resp.Close()

	buf, err := ioutil.ReadAll(r)
}
