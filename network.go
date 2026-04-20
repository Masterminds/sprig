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

func joinHostPort(host, port string) string {
	return net.JoinHostPort(host, port)
}

func splitHostPort(hostport string) map[string]string {
	host, port, _ := net.SplitHostPort(hostport)
	return map[string]string{"host": host, "port": port}
}
