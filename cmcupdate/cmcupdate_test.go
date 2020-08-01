package cmcupdate

import (
	"os"
	"path/filepath"
	"testing"

	Log "github.com/zytzjx/anthenacmc/loggersys"
)

func TestSaveStatusFile(t *testing.T) {
	saveStatusFile(([]byte)("aabbbb"))
}

func TestUpdateCMC(t *testing.T) {
	Log.NewLogger("updatecmc")
	UpdateCMC()
}

func TestDownloadCMC(t *testing.T) {
	Log.NewLogger("updatecmc")
	DownloadCMC("/opt/futuredial/hydradownloader")
}

func TestListFils(t *testing.T) {
	Log.NewLogger("updatecmc")
	files := map[string]bool{}
	err := filepath.Walk("/opt/futuredial/hydradownloader", visit(&files))
	if err != nil {
		Log.Log.Error(err)
	}
}

func TestRemoveAll(t *testing.T) {
	os.RemoveAll("/opt/futuredial/hydradownloader/*")
}
