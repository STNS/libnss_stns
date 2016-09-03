package libstns

import "net"

func NicReady() bool {
	is, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, i := range is {
		if i.Name[0:2] == "lo" || i.Flags&(1<<uint(0)) == 0 {
			continue
		}
		return true
	}
	return false
}
