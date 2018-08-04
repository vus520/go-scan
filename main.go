package main

import (
    "flag"
    "fmt"
    "github.com/the-404/IPgo"
    "net"
    "os"
    "strconv"
    "sync"
    "time"
)

var port int
var parallelCounts int
var verbose int
var timeOut int64
var protocol string

func init() {
    flag.IntVar(&port, "p", 5555, "检测的端口")
    flag.IntVar(&parallelCounts, "n", 10, "线程数")
    flag.IntVar(&verbose, "v", 0, "打印进度")
    flag.Int64Var(&timeOut, "t", 1500, "检测超时时间，毫秒")
    flag.StringVar(&protocol, "c", "tcp", "检测协议，tcp 或者 udp")

    // 修改提示信息
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "简单的TCP端口扫描器 pingo\nUsage: %s -[pnvt] <ip地址/网段>\n\nOptions:\n\n", "pingo")
        flag.PrintDefaults()
    }
    flag.Parse()
}

func printOpeningPort(ip net.IP, port int) {
    fmt.Println(ip.String() + ":" + strconv.Itoa(port) + " opening")
}

func checkPort(ip net.IP, port int, wg *sync.WaitGroup, parallelChan *chan int) {
    defer wg.Done()

    target := ip.String() + ":" + strconv.Itoa(port)

    if verbose > 0 {
        fmt.Println("checking " + target)
    }

    conn, err := net.DialTimeout(protocol, target, time.Millisecond*time.Duration(timeOut))
    if err == nil {
        printOpeningPort(ip, port)
        conn.Close()
    }
    <-*parallelChan
}

func main() {
    args := flag.Args()

    if len(args) != 1 {
        flag.Usage()

        return
    }

    if protocol != "udp" {
        protocol = "tcp"
    }

    ips := IPgo.Iplist(flag.Arg(0))

    if len(ips) == 1 && ips[0] == "<nil>" {
        fmt.Println("目标IP无法识别，请输入 IP地址 或 CIDR掩码\n\n")
        flag.Usage()
        return
    }

    // 用于协程任务控制
    wg := sync.WaitGroup{}
    parallelChan := make(chan int, parallelCounts)
    for i := 0; i < len(ips); i++ {
        ip := net.ParseIP(ips[i])
        wg.Add(1)
        parallelChan <- 1
        go checkPort(ip, port, &wg, &parallelChan)
    }
    wg.Wait()

}
