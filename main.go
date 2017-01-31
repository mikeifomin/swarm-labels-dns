package main

import (
	"./digitalocean"
	"./docker"
	"fmt"
	"github.com/joho/godotenv"
	"time"
)

func main() {
	for {
		fmt.Println("check")
		godotenv.Load()
		d := docker.NewFromEnv()
		do := digitalocean.NewFromEnv()
		nodesAddrs := d.FetchNodesAddrs()
		dnsAddrs := do.FetchDnsAddrs()

		for addr, nodeID := range nodesAddrs {
			domains := dnsAddrs[addr]
			fmt.Println(domains)
			d.UpdateNodeLabels(nodeID, domains)
		}
		time.Sleep(60 * time.Second)
	}
}
