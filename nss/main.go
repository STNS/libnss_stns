package main

import (
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

import "C"

func main() {
}

var pwdNss, grpNss, spwdNss *libstns.Nss

var pwdList = stns.Attributes{}
var pwdReadPos int

var grpList = stns.Attributes{}
var grpReadPos int

var spwdList = stns.Attributes{}
var spwdReadPos int

func init() {
	if pwdNss == nil && grpNss == nil && spwdNss == nil && libstns.AfterOsBoot() {
		libstns.Setlog()
		config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
		if err != nil {
			return
		}
		pwdNss = libstns.NewNss(config, "user", pwdList, &pwdReadPos)
		grpNss = libstns.NewNss(config, "group", grpList, &grpReadPos)
		spwdNss = libstns.NewNss(config, "user", spwdList, &spwdReadPos)
	}
}

func set(n *libstns.Nss, e libstns.NssEntry, column, value string) C.int {
	if !libstns.AfterOsBoot() || n == nil {
		return C.int(libstns.NSS_STATUS_NOTFOUND)
	}

	return C.int(n.Set(e, column, value))
}

func setByList(n *libstns.Nss, e libstns.NssEntry) C.int {
	if !libstns.AfterOsBoot() || n == nil {
		return C.int(libstns.NSS_STATUS_NOTFOUND)
	}

	return C.int(n.SetByList(e))
}

func initList(n *libstns.Nss, mode int) C.int {
	if !libstns.AfterOsBoot() || n == nil {
		return C.int(libstns.NSS_STATUS_NOTFOUND)
	}

	switch mode {
	case libstns.NSS_LIST_PRESET:
		return C.int(n.PresetList())
	case libstns.NSS_LIST_PURGE:
		n.PurgeList()
		return C.int(libstns.NSS_STATUS_SUCCESS)
	}
	return libstns.NSS_STATUS_NOTFOUND
}
