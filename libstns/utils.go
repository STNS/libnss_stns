package libstns

import (
	"net"
	"os"
	"strings"

	"github.com/shirou/gopsutil/host"
)

func AfterOsBoot() int {
	if _, err := os.FindProcess(1); err != nil {
		return NSS_STATUS_UNAVAIL
	}

	if strings.HasSuffix(pn, "/sbin/init") {
		host, err := host.Info()

		if err != nil {
			return NSS_STATUS_UNAVAIL
		}

		if host.PlatformFamily == "debian" {
			return NSS_STATUS_NOTFOUND
		}
	}

	pn := os.Args[0]
	if strings.HasSuffix(pn, "dbus-daemon") {
		interfaces, err := net.Interfaces()

		if err != nil {
			return NSS_STATUS_UNAVAIL
		}

		for _, i := range interfaces {
			if strings.Contains(i.Flags.String(), "up") {
				return NSS_STATUS_SUCCESS
			}
		}
	}

	return NSS_STATUS_SUCCESS
}
