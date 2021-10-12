package backup

import (
	"testing"
)

func Test_Zip(t *testing.T) {
	if err := Zip("backup.zip", "adir", "bdir"); err != nil {
		t.Error(err)
	}
}
