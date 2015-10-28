package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	err := StartLog("", DEBUG)
	if err != nil {
		t.Fatal(err)
	}
	Debug("test=======")
}
