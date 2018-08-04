pingo, 多线程端口扫描器
======================


### install

```shell

git clone https://github.com/vus520/go-scan.git
go get
go build -o pingo

./pingo

Usage: pingo -[prnv] <ip地址/网段>

Options:

  -n int
    	线程数 (default 10)
  -p int
    	检测的端口 (default 5555)
  -v int
    	打印进度

```


### run

```

pingo -p 80 -n 10 -v 1 192.168.50.1/24
checking 192.168.50.9:80
checking 192.168.50.4:80
checking 192.168.50.1:80
checking 192.168.50.10:80
192.168.50.1:80 opening

```

### more exsample
```
pingo -p 80 -n 10 -v 1 192.168.50.1
pingo -t 1000 192.168.50.1
```