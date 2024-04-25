package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IVEmmaWatson/scanner"
)

var (
	DstIP string
	GW    string
	SrcIP string
	port  string
	help  bool
)

func main() {

	flag.StringVar(&DstIP, "d", "111.111.111.111", "目标IP地址")                               // 有检查
	flag.StringVar(&GW, "g", "192.168.0.1", "本地网关IP地址")                                  // 用户输错网关ip，用户承担风险
	flag.StringVar(&SrcIP, "s", "192.168.0.2/24", "本机IP地址,需为CIDR表示法如192.168.0.2/24") // 有检查
	flag.StringVar(&port, "p", "79-81", "目标端口范围")                                        //  有检查
	flag.BoolVar(&help, "help", false, "Usage: show help message")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	ipResult, err := ppp.ParseIP(DstIP)
	portResult, err := ppp.ParsePort(port)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for _, ip := range ipResult {
		err := ppp.SynScan(ip, GW, SrcIP, portResult)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

}
