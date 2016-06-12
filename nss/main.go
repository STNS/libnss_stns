package main

import (
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)
import "C"

func main() {
}

func init() {
	libstns.Setlog()
}

func set(s libstns.SetNss, entry, presult interface{}, r, c, v string) C.int {
	nss, err := libstns.NewNss(r, c, v)
	if err != nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
	}
	return C.int(nss.Set(s, entry, presult))
}

func setByList(s libstns.SetNss, entry, presult interface{}, list stns.Attributes, position *int) C.int {
	nss, err := libstns.NewNss("", "", "")
	if err != nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
	}
	return C.int(nss.SetByList(s, entry, presult, list, position))
}

func initList(list stns.Attributes, position *int, r string, mode int) C.int {
	nss, err := libstns.NewNss(r, "list", "")
	if err != nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
	}

	switch mode {
	case libstns.NSS_LIST_PRESET:
		return C.int(nss.PresetList(list, position))
	case libstns.NSS_LIST_PURGE:
		nss.PurgeList(list, position)
		return C.int(libstns.NSS_STATUS_SUCCESS)
	}
	return libstns.NSS_STATUS_NOTFOUND
}
