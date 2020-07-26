package reportcmc

import (
	"fmt"
	"testing"

	Log "github.com/zytzjx/anthenacmc/loggersys"
)

func TestNewUUID(t *testing.T) {
	uuid, err := newUUID()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("%s\n", uuid)
}

func BenchmarkNewUUID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uuid, err := newUUID()
		if err != nil {

		}
		fmt.Printf("%s\n", uuid)
	}
}

func TestReportCMC(t *testing.T) {
	Log.NewLogger("reportcmc")
	ReportCMC()
}
