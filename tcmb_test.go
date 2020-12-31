package tcmb

import (
	"testing"
)

func Test_Get(t *testing.T) {
	res, err := Get()
	if err != nil {
		t.Error("Failed to get.")
		return
	}

	if len(res.Currencies) == 0 {
		t.Error("length of currencies is zero")
		return
	}

}

