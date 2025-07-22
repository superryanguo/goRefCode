package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
)

var (
	addr = flag.String("addr", ":8080", "service address")
)

func main() {
	flag.Parse()
	lc := &net.ListenConfig{}
	if lc.MultipathTCP() { // 默认mptcp是禁用的
		panic("MultipathTCP should be off by default")
	}
	lc.SetMultipathTCP(true)                                 // 主动启用mptcp
	ln, err := lc.Listen(context.Background(), "tcp", *addr) // 正常tcp监听
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			isMultipathTCP, err := conn.(*net.TCPConn).MultipathTCP() // 检查连接是否支持了mptcp
			fmt.Printf("accepted connection from %s with mptcp: %t, err: %v\n", conn.RemoteAddr(), isMultipathTCP, err)
			for {
				clientIP, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					panic(err)
				}
				dataStr := string(buf[:n])
				fmt.Printf("clientIP: %s  data: %s \n", clientIP, dataStr)
			}
		}()
	}
}
