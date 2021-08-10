package reportcmc

import (
	"encoding/json"
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
	ReportCMC("")
}

func TestReportResult(t *testing.T) {
	Log.NewLogger("reportcmc")
	msg := `{"recordid": "12808398b52f4c53b7c328dea24ee899", "status": 1}`
	var status map[string]interface{}
	// fmt.Println(string(resp.Body()))
	if err := json.Unmarshal([]byte(msg), &status); err != nil {
		// fmt.Println("return format error.")
		Log.Log.Error("return format error.")
		// return 0, err
	}

	if val, ok := status["status"]; ok {
		fmt.Printf("(%v, %T)\n", val, val)
		if s, ok := val.(string); ok {
			switch s {
			case "1":
				// return 1, nil
				fmt.Println("1")
			case "2":
				// fmt.Println("success")
				Log.Log.Info("Success")
				fmt.Println("2")
			case "3":
				// return 3, fmt.Errorf("%v", status["error"])
				fmt.Println("3")
			case "4":
				Log.Log.Warn("same uuid")
				// fmt.Println("same uuid")
				// return 4, nil
				fmt.Println("4")
			default:
				// fmt.Println("failed.")
				Log.Log.Error("failed.")
			}
		} else if vv, ok := val.(float64); ok {
			stat := int(vv)
			switch stat {
			case 1:
				// return 1, nil
				fmt.Println("1")
			case 2:
				// fmt.Println("success")
				Log.Log.Info("Success")
				// return 2, nil
				fmt.Println("2")
			case 3:
				// return 3, fmt.Errorf("%v", status["error"])
				fmt.Println("3")
			case 4:
				Log.Log.Warn("same uuid")
				// fmt.Println("same uuid")
				// return 4, nil
				fmt.Println("4")
			default:
				// fmt.Println("failed.")
				Log.Log.Error("failed.")
			}
		}

	}
	// return 0, errors.New("I do not know the protocol return")
}
