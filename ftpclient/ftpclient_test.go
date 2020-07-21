package ftpclient

import "testing"

func TestCreateFolder(t *testing.T) {
	if err := CreateFolder("abc"); err != nil {
		t.Error(err)
	}
}
