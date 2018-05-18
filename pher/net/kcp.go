package net

import (
	"fmt"
	"net"

	"github.com/fatedier/frp/utils/log"

	kcp "github.com/fatedier/kcp-go"
)

type KcpListener struct {
	net.Addr
	listener  net.Listener
	accept    chan Conn
	closeFlag bool
	log.Logger
}

func ListenKcp(bindAddr string, bindPort int) (l *KcpListener, err error) {
	listener, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", bindAddr, bindPort), nil, 10, 3)
	if err != nil {
		return l, err
	}
	listener.SetReadBuffer(4194304)
	listener.SetWriteBuffer(4194304)

	l = &KcpListener{
		Addr:      listener.Addr(),
		listener:  listener,
		accept:    make(chan Conn),
		closeFlag: false,
		Logger:    log.NewPrefixLogger(""),
	}

	go func() {
		for {
			conn, err := listener.AcceptKCP()
			if err != nil {
				if l.closeFlag {
					close(l.accept)
					return
				}
				continue
			}
			conn.SetStreamMode(true)
			conn.SetWriteDelay(true)
			conn.SetNoDelay(1, 20, 2, 1)
			conn.SetMtu(1350)
			conn.SetWindowSize(1024, 1024)
			conn.SetACKNoDelay(false)

			l.accept <- WrapConn(conn)
		}
	}()
	return l, err
}

func (l *KcpListener) Accept() (Conn, error) {
	conn, ok := <-l.accept
	if !ok {
		return conn, fmt.Errorf("channel for kcp listener closed")
	}
	return conn, nil
}

func (l *KcpListener) Close() error {
	if !l.closeFlag {
		l.closeFlag = true
		l.listener.Close()
	}
	return nil
}

func NewKcpConnFromUdp(conn *net.UDPConn, connected bool, raddr string) (net.Conn, error) {
	kcpConn, err := kcp.NewConnEx(1, connected, raddr, nil, 10, 3, conn)
	if err != nil {
		return nil, err
	}
	kcpConn.SetStreamMode(true)
	kcpConn.SetWriteDelay(true)
	kcpConn.SetNoDelay(1, 20, 2, 1)
	kcpConn.SetMtu(1350)
	kcpConn.SetWindowSize(1024, 1024)
	kcpConn.SetACKNoDelay(false)
	return kcpConn, nil
}
