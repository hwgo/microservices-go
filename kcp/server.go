package main

import (
	"fmt"

	"github.com/fatedier/frp/utils/log"
	frpNet "github.com/fatedier/frp/utils/net"
)

var (
	bindAddr    = "127.0.0.1"
	kcpBindPort = 10001
)

// Server service.
type Service struct {
	// Accept connections from client.
	listener frpNet.Listener

	// Accept connections using kcp.
	kcpListener frpNet.Listener
}

func NewServer() (svr *Service, err error) {

	svr = &Service{}

	svr.listener, err = frpNet.ListenKcp(bindAddr, kcpBindPort)
	if err != nil {
		err = fmt.Errorf("Listen on kcp address udp [%s:%d] error: %v", bindAddr, kcpBindPort, err)
		return
	}
	log.Info("frps kcp listen on udp %s:%d", bindAddr, kcpBindPort)

	for {
		_, err := l.Accept()
		if err != nil {
			log.Warn("Listener for incoming connections from client closed")
			return
		}
	}
}

func (svr *Service) Run() {
	svr.HandleListener(svr.kcpListener)
}

func (svr *Service) HandleListener(l frpNet.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Warn("Listener for incoming connections from client closed")
			return
		}

		// Start a new goroutine for dealing connections.
		go func(frpConn frpNet.Conn) {
			dealFn := func(conn frpNet.Conn) {
				var rawMsg msg.Message
				conn.SetReadDeadline(time.Now().Add(connReadTimeout))
				if rawMsg, err = msg.ReadMsg(conn); err != nil {
					log.Trace("Failed to read message: %v", err)
					conn.Close()
					return
				}
				conn.SetReadDeadline(time.Time{})

				switch m := rawMsg.(type) {
				case *msg.Login:
					err = svr.RegisterControl(conn, m)
					// If login failed, send error message there.
					// Otherwise send success message in control's work goroutine.
					if err != nil {
						conn.Warn("%v", err)
						msg.WriteMsg(conn, &msg.LoginResp{
							Version: version.Full(),
							Error:   err.Error(),
						})
						conn.Close()
					}
				case *msg.NewWorkConn:
					svr.RegisterWorkConn(conn, m)
				case *msg.NewVisitorConn:
					if err = svr.RegisterVisitorConn(conn, m); err != nil {
						conn.Warn("%v", err)
						msg.WriteMsg(conn, &msg.NewVisitorConnResp{
							ProxyName: m.ProxyName,
							Error:     err.Error(),
						})
						conn.Close()
					} else {
						msg.WriteMsg(conn, &msg.NewVisitorConnResp{
							ProxyName: m.ProxyName,
							Error:     "",
						})
					}
				default:
					log.Warn("Error message type for the new connection [%s]", conn.RemoteAddr().String())
					conn.Close()
				}
			}

			if config.ServerCommonCfg.TcpMux {
				session, err := smux.Server(frpConn, nil)
				if err != nil {
					log.Warn("Failed to create mux connection: %v", err)
					frpConn.Close()
					return
				}

				for {
					stream, err := session.AcceptStream()
					if err != nil {
						log.Warn("Accept new mux stream error: %v", err)
						session.Close()
						return
					}
					wrapConn := frpNet.WrapConn(stream)
					go dealFn(wrapConn)
				}
			} else {
				dealFn(frpConn)
			}
		}(c)
	}
}

func main() {
	fmt.Println(1)
}
