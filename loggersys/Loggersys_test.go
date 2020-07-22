package loggersys

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	loger := NewLogger("aa")
	loger.Info("Ad ddd ddd ddd ")
}
