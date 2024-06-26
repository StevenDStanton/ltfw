package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/StevenDStanton/cli-tools-for-windows/common"
	"github.com/StevenDStanton/cli-tools-for-windows/crypto/api"
)

var (
	args []string
)

const (
	tool    = "Crypto"
	version = "v1.0.7"
)

func init() {
	args = os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("Must specify at least one pair such as BTC/USD")
	}
	versionInformation := common.PrintVersion(tool, version)
	fmt.Println(versionInformation)

}

func fetchRate(pair string) {
	response, err := api.GetRate(pair)
	if err != nil {
		fmt.Printf("Error fetching rate for %s: %v", pair, err)
		return
	}
	fmt.Println(response)
}

func main() {
	var wg sync.WaitGroup
	for _, pair := range args {
		wg.Add(1)
		go func() {
			fetchRate(pair)
			wg.Done()
		}()
	}
	wg.Wait()
}
