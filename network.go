package sprig

import (
	"math/rand"
	"net"
)

func getHostByName(name string) string {
	addrs, err := net.LookupHost(name)
	if err != nil || len(addrs) == 0 {
		return name
	}
	return addrs[rand.Intn(len(addrs))]
}
