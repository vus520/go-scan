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

func init() {
    flag.IntVar(&port, "p", 5555, "检测的端口")
    flag.IntVar(&parallelCounts, "n", 10, "线程数")
    flag.IntVar(&verbose, "v", 0, "打印进度")

    // 修改提示信息
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "简单的TCP端口扫描器 pingo\nUsage: %s -[prn] <ip地址/网段>\n\nOptions:\n\n", "pingo")
        flag.PrintDefaults()
    }
    flag.Parse()
}

func printOpeningPort(ip net.IP, port int) {
    fmt.Println(ip.String() + ":" + strconv.Itoa(port) + " opening")
}

func checkPort(ip net.IP, port int, wg *sync.WaitGroup, parallelChan *chan int) {

    target := ip.String() + ":" + strconv.Itoa(port)

    if verbose > 0 {
        fmt.Println("checking " + target)
    }

    defer wg.Done()

    conn, err := net.DialTimeout("tcp", target, time.Second*1)
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

    ips := IPgo.Iplist(flag.Arg(0))

    if len(ips) == 1 && ips[0] == "<nil>" {
        fmt.Println("请输入 IP地址 或 CIDR掩码\n\n")
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
