package cache

import (
	"errors"
	"testing"

	"github.com/STNS/STNS/stns"
)

func TestWriteRead(t *testing.T) {
	SetWorkDir("/tmp")

	Write("test", stns.Attributes{"test": &stns.Attribute{
		Id: 1,
	}},
		nil,
	)

	attrs, err := Read("test")

	if attrs["test"].Id != 1 || err != nil {
		t.Error("rw error1")
	}

	Write("error", stns.Attributes{"error": &stns.Attribute{
		Id: 1,
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
		Id: 1,
	}},
	)

	attrs := *LastResultList("test")
	if attrs["test"].Id != 1 {
		t.Error("save load error1")
	}
}
