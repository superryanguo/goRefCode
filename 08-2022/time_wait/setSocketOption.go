package main

import (
	"context"
	"net"
	"net/http"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	//我忽然想起来以前在《UNIX 网络编程》上有看到过一个Socket的参数，叫 <code>SO_LINGER，我的编程生涯中从来没有使用过这个设置，这个参数主要是为了延尽关闭来用的，也就是说你应用调用 close()函数时，如果还有数据没有发送完成，则需要等一个延时时间来让数据发完，但是，如果你把延时设置为 0  时，Socket就丢弃数据，并向对方发送一个 RST 来终止连接，因为走的是 RST 包，所以就不会有 TIME_WAIT 了。

	//这个东西在服务器端永远不要设置，不然，你的客户端就总是看到 TCP 链接错误 “connnection reset by peer”，但是这个参数对于 EaseProbe 的客户来说，简直是太完美了，当EaseProbe 探测完后，直接 reset connection， 即不会有功能上的问题，也不会影响服务器，更不会有烦人的 TIME_WAIT 问题。
	//在 Golang的标准库代码里，net.TCPConn 有个方法 SetLinger()可以完成这个事，使用起来也比较简单：
	conn, _ := net.DialTimeout("tcp", t.Host, t.Timeout())

	if tcpCon, ok := conn.(*net.TCPConn); ok {
		tcpCon.SetLinger(0)
	}
	//++++++++
	dialer := &net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				operr = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.TCP_QUICKACK, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}
	//上面这个方法非常的低层，需要直接使用setsocketopt这样的系统调用，我其实，还是想使用 TCPConn.SetLinger(0) 来完成这个事，我不是很难碰底层的事。
	//经过Golang http包的源码阅读和摸索，我使用了下面的方法
	client := &http.Client{
		Timeout: h.Timeout(),
		Transport: &http.Transport{
			TLSClientConfig:   tls,
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{Timeout: h.Timeout()}
				conn, err := d.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				tcpConn, ok := conn.(*net.TCPConn)
				if ok {
					tcpConn.SetLinger(0)
					return tcpConn, nil
				}
				return conn, nil
			},
		},
	}
}
