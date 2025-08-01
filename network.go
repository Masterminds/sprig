package sprig

import (
	"net"
)

const NUM_TRIES = 3

func getHostByName(name string) ([]string, error) {
	err := error(nil)
	addrs := []string(nil)
	for tries := 0; tries < NUM_TRIES; tries++ {
		if addrs, err = net.LookupHost(name); err == nil {
			return addrs, nil
		}
	}
	return addrs, err
}
