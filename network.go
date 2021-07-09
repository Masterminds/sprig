package sprig

import (
	"math/rand"
	"net"
)

func getHostByName(name string) string {
	addrs, err := net.LookupHost(name)
	if err != nil {
		return ""
	}

	return addrs[rand.Intn(len(addrs))]
}
