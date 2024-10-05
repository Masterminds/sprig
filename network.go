package sprig

import (
	"math/rand"
	"net"
	"fmt"
	"time"
)

// Error check loop runs 3 times with 15 second intervals only if DNS/connection type errors,
// if host does not exist "" blank gets returned
func getHostByName(name string) (string, error) {
	err := error(nil)
	for tries := 3; tries > 0; tries-- {
		addrs, err := net.LookupHost(name)
		if err == nil {
			return addrs[rand.Intn(len(addrs))], nil
		}
		dnsErr, ok := err.(*net.DNSError)
		if !ok {
			return "", fmt.Errorf("DNS Error, %v", err)
		}
		if dnsErr.IsNotFound {
			return "", nil
		}
		time.Sleep(15 * time.Second)
	}
	return "", fmt.Errorf("Failure in looking up %s: ", err)
}
