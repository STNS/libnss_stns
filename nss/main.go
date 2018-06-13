package main

import (
	"os"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
	"github.com/shirou/gopsutil/host"
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
	orgInit()
}

func ignoreProcess(process string) bool {
	ignore := []string{
		"/sbin/init",
		"dbus-daemon",
		"(resolved)",
		"(systemd)",
	}

	for _, i := range ignore {
		if strings.HasSuffix(process, i) {
			return true
		}
	}
	return false
}

func orgInit() int {
	if pwdNss == nil || grpNss == nil || spwdNss == nil {
		if _, err := os.FindProcess(1); err != nil {
			return libstns.NSS_STATUS_UNAVAIL
		}

		host, err := host.Info()
		if err != nil {
			return libstns.NSS_STATUS_UNAVAIL
		}

		if host.PlatformFamily == "debian" && ignoreProcess(os.Args[0]) {
			return libstns.NSS_STATUS_NOTFOUND
		}
		libstns.Setlog()

		config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
		if err != nil {
			return libstns.NSS_STATUS_NOTFOUND
		}

		pwdNss = libstns.NewNss(config, "user", pwdList, &pwdReadPos)
		grpNss = libstns.NewNss(config, "group", grpList, &grpReadPos)
		spwdNss = libstns.NewNss(config, "user", spwdList, &spwdReadPos)
	}
	return libstns.NSS_STATUS_SUCCESS
}

func set(n *libstns.Nss, e libstns.NssEntry, column, value string) C.int {
	a := orgInit()
	if a != libstns.NSS_STATUS_SUCCESS {
		return C.int(a)
	}

	if n == nil {
		return C.int(libstns.NSS_STATUS_NOTFOUND)
	}

	return C.int(n.Set(e, column, value))
}

func setByList(n *libstns.Nss, e libstns.NssEntry) C.int {
	a := orgInit()
	if a != libstns.NSS_STATUS_SUCCESS {
		return C.int(a)
	}

	if n == nil {
		return C.int(libstns.NSS_STATUS_NOTFOUND)
	}

	return C.int(n.SetByList(e))
}

func initList(n *libstns.Nss, mode int) C.int {
	a := orgInit()
	if a != libstns.NSS_STATUS_SUCCESS {
		return C.int(a)
	}

	if n == nil {
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
