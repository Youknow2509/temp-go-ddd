package grpc

import (
	"net"
	"time"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// tcpTunedListener wraps *net.TCPListener to apply TCP options on Accepted connections.
type tcpTunedListener struct {
	*net.TCPListener
	cfg domain_config.GrpcTcpSetting
}

func newTCPTunedListener(ln *net.TCPListener, tcpCfg domain_config.GrpcTcpSetting) *tcpTunedListener {
	return &tcpTunedListener{TCPListener: ln, cfg: tcpCfg}
}

func (l *tcpTunedListener) Accept() (net.Conn, error) {
	conn, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}

	_ = conn.SetNoDelay(l.cfg.TcpNodelay)

	if l.cfg.TcpKeepalive {
		_ = conn.SetKeepAlive(true)
		if l.cfg.TcpKeepaliveTimeMs > 0 {
			_ = conn.SetKeepAlivePeriod(time.Duration(l.cfg.TcpKeepaliveTimeMs) * time.Millisecond)
		}
	} else {
		_ = conn.SetKeepAlive(false)
	}

	return conn, nil
}
