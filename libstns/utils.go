package libstns

import (
	"net"
	"strings"
)

func NicReady() bool {
	i, err := net.InterfaceByIndex(0)
	if err != nil || !strings.Contains(i.Flags.String(), "up") {
		return false
	}
	return true
}
