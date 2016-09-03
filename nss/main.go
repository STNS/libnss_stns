package main

import (
	"log"

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
	libstns.Setlog()

	if !libstns.NicReady() {
		log.Println("does not have a valid network interface")
		return
	}

	config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
	if err != nil {
		return
	}

	var e error
	pwdNss, e = libstns.NewNss(config, "user", pwdList, &pwdReadPos)
	if e != nil {
		return
	}

	grpNss, e = libstns.NewNss(config, "group", grpList, &grpReadPos)
	if e != nil {
		return
	}

	spwdNss, e = libstns.NewNss(config, "user", spwdList, &spwdReadPos)
	if e != nil {
		return
	}
}

func set(n *libstns.Nss, e libstns.NssEntry, column, value string) C.int {
	if n == nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
	}
	return C.int(n.Set(e, column, value))
}

func setByList(n *libstns.Nss, e libstns.NssEntry) C.int {
	if n == nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
	}
	return C.int(n.SetByList(e))
}

func initList(n *libstns.Nss, mode int) C.int {
	if n == nil {
		return C.int(libstns.NSS_STATUS_UNAVAIL)
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
