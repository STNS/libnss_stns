package libstns

import (
	"os"

	"github.com/shirou/gopsutil/host"
)

func AfterOsBoot() bool {
	_, err := os.FindProcess(1)
	if err != nil {
		host, err := host.Info()
		if err != nil || (host.PlatformFamily == "debian" && (os.Args[0] == "/sbin/init" || os.Args[0] == "dbus-daemon")) {
			return false
		}
	}
	return true
}
