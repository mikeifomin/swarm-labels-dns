package main

import (
	"./digitalocean"
	"./docker"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("dsds")
	godotenv.Load()
	d := docker.NewFromEnv()
	do := digitalocean.NewFromEnv()
	nodesAddrs := d.FetchNodesAddrs()
	dnsAddrs := do.FetchDnsAddrs()

	for addr, nodeID := range nodesAddrs {
		domains := dnsAddrs[addr]
		d.UpdateNodeLabels(nodeID, domains)
	}
}
