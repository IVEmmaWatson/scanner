package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IVEmmaWatson/scanner"
)

var (
	dstIP string
	gw    string
	srcIP string
	port  string
	help  bool
)

func init() {
	flag.StringVar(&dstIP, "d", "111.111.111.111", "目标IP地址")
	flag.StringVar(&gw, "g", "192.168.0.1", "本地网关IP地址")
	flag.StringVar(&srcIP, "s", "192.168.0.2/24", "本机IP地址,需为CIDR表示法如192.168.0.2/24")
	flag.StringVar(&port, "p", "79-81", "目标端口范围")
	flag.BoolVar(&help, "help", false, "Usage: show help message")
}

func main() {
	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	ipResult, err := scanner.ParseIP(dstIP)
	portResult, err := scanner.ParsePort(port)
	if err != nil {
		log.Fatalf("无法解析目标地址: %v", err)
	}

	for _, ip := range ipResult {
		fmt.Println("--------------------------")
		fmt.Printf("scanner ip:%v\n", ip)
		err := scanner.SynScan(ip, gw, srcIP, portResult)
		if err != nil {
			log.Printf("扫描失败: %v", err)
		}
	}

}
