package loggersys

import(
	"testing"
)

func TestNewLogger(t *testing.T)  {
	loger:=NewLogger() 
	loger.Info("Ad ddd ddd ddd ")
}