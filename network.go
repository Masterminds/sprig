package sprig

import (
	"fmt"
	"math/rand"
	"net"
)

func getHostByName(name string) (string, error) {
	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("unable to resolve %s: %v", name, err)
		return "", err
	}
	return addrs[rand.Intn(len(addrs))], nil
}
