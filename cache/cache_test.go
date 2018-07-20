package cache

import (
	"errors"
	"os"
	"testing"

	"github.com/STNS/STNS/stns"
)

func TestWriteRead(t *testing.T) {
	SetWorkDir("/tmp")

	Write("test", stns.Attributes{"test": &stns.Attribute{
		ID: 1,
	}},
		nil,
	)

	attrs, err := Read("test")

	if attrs["test"].ID != 1 || err != nil {
		t.Error("rw error1")
	}

	Write("error", stns.Attributes{"error": &stns.Attribute{
		ID: 1,
	}},
		errors.New("test error"),
	)

	attrs, err = Read("error")
	if attrs != nil || err.Error() != "test error" {
		t.Error("rw error2")
	}
}

func TestSaveResultList(t *testing.T) {
	SetWorkDir("/tmp")

	SaveResultList("test", stns.Attributes{"test": &stns.Attribute{
		ID: 1,
	}},
	)

	cachePath := "/tmp/.libnss_stns_test_cache"
	stat, err := os.Stat(cachePath)

	if err != nil {
		t.Errorf("Got error %v", err)
	}
	if mode := stat.Mode(); mode != 0660 {
		t.Errorf("Expect 0660, got %v", mode)
	}

}

func TestSaveLoad(t *testing.T) {
	SetWorkDir("/tmp")

	SaveResultList("test", stns.Attributes{"test": &stns.Attribute{
		ID: 1,
	}},
	)

	attrs := *LastResultList("test")
	if attrs["test"].ID != 1 {
		t.Error("save load error1")
	}
}

func TestWriteReadMinID(t *testing.T) {
	WriteID("test", "min", 100)
	if ReadMinID("test") != 100 {
		t.Error("minid error1")
	}
}

func TestWriteReadMaxID(t *testing.T) {
	WriteID("test", "max", 100)
	if ReadMaxID("test") != 100 {
		t.Error("max id error1")
	}
}
