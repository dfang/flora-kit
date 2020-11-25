package flora

import (
	"errors"
	"net"
	"sync"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type ShadowSocksServer struct {
	proxyType string
	server    string
	// cipher    *ss.Cipher
	cipher    core.Cipher
	failCount int
	lock      sync.RWMutex
}

func NewShadowSocks(server string, cipher core.Cipher) *ShadowSocksServer {
	return &ShadowSocksServer{
		proxyType: ServerTypeShadowSocks,
		server:    server,
		cipher:    cipher,
	}
}

func (s *ShadowSocksServer) ResetFailCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.failCount = 0
}

func (s *ShadowSocksServer) AddFail() {
	s.failCount++
}

func (s *ShadowSocksServer) FailCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.failCount
}

func (s *ShadowSocksServer) ProxyType() string {
	return s.proxyType
}

func (s *ShadowSocksServer) DialWithRawAddr(raw []byte, host string) (net.Conn, error) {
	// if nil != raw && len(raw) > 0 {
	// 	return ss.DialWithRawAddr(raw, s.server, s.cipher.Copy())
	// } else {
	// 	return ss.Dial(host, s.server, s.cipher.Copy())
	// }

	// fmt.Println("Host is", host)
	// fmt.Println("s.server is", s.server)
	conn, err := core.Dial("tcp", s.server, s.cipher) //core code
	if err != nil {
		return nil, err
	}
	tgt := socks.ParseAddr(host)
	if tgt == nil {
		return nil, errors.New("wrong dest address")
	}
	conn.Write(tgt)

	return conn, nil
}
