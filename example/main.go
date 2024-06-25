package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/IVEmmaWatson/scanner"
)

var (
	dstIP     string
	gw        string
	srcIP     string
	port      string
	help      bool
	numWorker int
)

func init() {
	flag.StringVar(&dstIP, "d", "47.120.52.232", "目标IP地址")
	flag.StringVar(&gw, "g", "192.168.3.1", "本地网关IP地址")
	flag.StringVar(&srcIP, "s", "192.168.3.19/24", "本机IP地址,需为CIDR表示法如192.168.0.2/24")
	flag.StringVar(&port, "p", "0-65535", "目标端口范围")
	flag.BoolVar(&help, "help", false, "Usage: show help message")
	flag.IntVar(&numWorker, "w", 3, "工作池大小") //
}

// 只读通道，读取任务，执行任务
func worker(jobs <-chan []uint16, wg *sync.WaitGroup, ipResult []net.IP, gw, srcIP string) {

	for job := range jobs {
		for _, i2 := range ipResult {
			fmt.Println("-----------")
			fmt.Printf("scanner ip:%v\n", i2)
			err := scanner.SynScan(i2, gw, srcIP, job)
			if err != nil {
				log.Printf("扫描失败: %v", err)
			}
		}
		wg.Done()
	}

}

func main() {

	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	start := time.Now()

	ipResult, err := scanner.ParseIP(dstIP)
	portResult, err := scanner.ParsePort(port)
	if err != nil {
		log.Fatalf("无法解析目标地址: %v", err)
	}

	// 创建任务池,以端口为任务
	jobs := make(chan []uint16)

	var wg sync.WaitGroup

	// 开启任务，最大同时工作量为numWorker
	for i := 0; i < numWorker; i++ {
		go worker(jobs, &wg, ipResult, gw, srcIP)
	}

	size := 10000

	// 分配任务信息
	for i := 0; i < len(portResult); i += size {
		end := i + size
		if end > len(portResult) {
			end = len(portResult)
		}
		wg.Add(1)
		jobs <- portResult[i:end]
	}

	// 任务清空后，关闭通道
	close(jobs)

	// 等待全部goroutine完成
	wg.Wait()
	ak := time.Since(start)
	fmt.Println("---扫描结束---")
	fmt.Println("总耗时:", ak)
}
