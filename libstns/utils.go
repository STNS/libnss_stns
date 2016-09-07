package libstns

import (
	"os"

	"github.com/shirou/gopsutil/host"
)

func AfterOsBoot() int {
	host, err := host.Info()

	if err != nil {
		return NSS_STATUS_UNAVAIL
	}

	if host.PlatformFamily == "debian" && (os.Args[0] == "/sbin/init" || os.Args[0] == "dbus-daemon") {
		return NSS_STATUS_NOTFOUND
	}

	return NSS_STATUS_SUCCESS
}
