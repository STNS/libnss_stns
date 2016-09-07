package libstns

import "os"

func AfterOsBoot() bool {
	// mount check
	_, err := os.FindProcess(1)
	if err != nil || os.Args[0] == "/sbin/init" || os.Args[0] == "dbus-daemon" {
		return false
	}
	return true
}
