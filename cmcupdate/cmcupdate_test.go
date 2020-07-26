package cmcupdate

import (
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
	DownloadCMC()
}
