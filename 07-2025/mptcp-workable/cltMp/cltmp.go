package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var (
	addr = flag.String("addr", "172.21.86.184:8080", "service address")


func main() {
	flag.Parse()
	d := &net.Dialer{}
	if d.MultipathTCP() { // 默认不启用
		panic("MultipathTCP should be off by default")
	}
	d.SetMultipathTCP(true) // 主动启用
	if !d.MultipathTCP() {  // 已经设置dial的时候使用mptcp
		panic("MultipathTCP is not on after having been forced to on")
	}
	c, err := d.Dial("tcp", *addr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	tcp, ok := c.(*net.TCPConn)
	if !ok {
		panic("struct is not a TCPConn")
	}
	mptcp, err := tcp.MultipathTCP() // 建立的连接是否真的支持mptcp
	if err != nil {
		panic(err)
	}
	fmt.Printf("outgoing connection from %s with mptcp: %t\n", *addr, mptcp)
	if !mptcp { // 不支持mptcp, panic
		panic("outgoing connection is not with MPTCP")
	}
	for {
		snt := []byte("MPTCP TEST")
		if _, err := c.Write(snt); err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
