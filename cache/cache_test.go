package cache

import (
	"errors"
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
	WriteMinID("test", 100)
	if ReadMinID("test") != 100 {
		t.Error("min id error1")
	}
}
