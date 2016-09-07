package libstns

import (
	"os"

	"github.com/shirou/gopsutil/host"
)

func AfterOsBoot() int {
	if _, err := os.FindProcess(1); err != nil {
		return NSS_STATUS_UNAVAIL
	}

	if os.Args[0] == "/sbin/init" || os.Args[0] == "dbus-daemon" {
		host, err := host.Info()

		if err != nil {
			return NSS_STATUS_UNAVAIL
		}

		if host.PlatformFamily == "debian" {
			return NSS_STATUS_NOTFOUND
		}
	}
	return NSS_STATUS_SUCCESS
}
