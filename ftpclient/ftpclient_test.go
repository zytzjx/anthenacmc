package ftpclient

import (
	"testing"

	Log "github.com/zytzjx/anthenacmc/loggersys"
)

func TestCreateFolder(t *testing.T) {
	fi := newFTPInfo()
	if err := fi.CreateFolder("abc"); err != nil {
		t.Error(err)
	}
}

func TestPrintInfo(t *testing.T) {
	Log.NewLogger("ftp")
	fi := newFTPInfo()
	fi.PrintInfo()
}
