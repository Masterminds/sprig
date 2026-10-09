package sprig

import (
	"fmt"
	"math/rand"
	"net"
)

func getHostByName(name string) string {
	addrs, _ := net.LookupHost(name)
	//TODO: add error handing when release v3 comes out
	return addrs[rand.Intn(len(addrs))]
}

func joinHostPort(host, port interface{}) string {
	return net.JoinHostPort(fmt.Sprint(host), fmt.Sprint(port))
}

