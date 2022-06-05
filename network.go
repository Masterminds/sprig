package sprig

import (
	"math/rand"
	"net"
)

func getHostByName(name string) string {
	addrs, _ := net.LookupHost(name)
	//TODO: add error handing when release v3 comes out
	return addrs[rand.Intn(len(addrs))]
}

func getHostByAddr(addr string) string {
	hosts, _ := net.LookupAddr(addr)
	return hosts[0]
}
