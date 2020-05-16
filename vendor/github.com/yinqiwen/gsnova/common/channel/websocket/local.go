package websocket

import (
	"crypto/tls"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/yinqiwen/gsnova/common/channel"
	"github.com/yinqiwen/gsnova/common/logger"
	"github.com/yinqiwen/gsnova/common/mux"
	"github.com/yinqiwen/pmux"
)

type WebsocketProxy struct {
}

func (p *WebsocketProxy) Features() channel.FeatureSet {
	return channel.FeatureSet{
		AutoExpire: true,
		Pingable:   true,
	}
}

func (ws *WebsocketProxy) CreateMuxSession(server string, conf *channel.ProxyChannelConfig) (mux.MuxSession, error) {
	u, err := url.Parse(server)
	if nil != err {
		return nil, err
	}
	u.Path = "/ws"
	wsDialer := &websocket.Dialer{}
	wsDialer.NetDial = channel.NewDialByConf(conf, u.Scheme)
	wsDialer.TLSClientConfig = channel.NewTLSConfig(conf)

	if conf.ForceTls13 == "tls13" {
		wsDialer.TLSClientConfig.MaxVersion = tls.VersionTLS13
		wsDialer.TLSClientConfig.MinVersion = tls.VersionTLS13
		wsDialer.TLSClientConfig.CipherSuites = []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		}
		logger.Info("wss ForceTls13, Now tls13 dial.")
	} else if conf.ForceTls13 == "tls12" {
		wsDialer.TLSClientConfig.MaxVersion = tls.VersionTLS12
		wsDialer.TLSClientConfig.MinVersion = tls.VersionTLS12
		wsDialer.TLSClientConfig.CipherSuites = []uint16{
			//tls.TLS_AES_128_GCM_SHA256,
			//tls.TLS_AES_256_GCM_SHA384,
			//tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		}
		logger.Info("wss ,tls12 dial.")
	} else { // "" or "auto" or other
		wsDialer.TLSClientConfig.MaxVersion = tls.VersionTLS13
		wsDialer.TLSClientConfig.MinVersion = tls.VersionTLS12
		wsDialer.TLSClientConfig.CipherSuites = []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		}
		logger.Info("wss ,tls12-tls13 auto dial.")
	}

	c, _, err := wsDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Notice("dial websocket error:%v %v", err, u.String())
		return nil, err
	}
	logger.Info("Connect %s success from %v->%v", server, c.LocalAddr(), c.RemoteAddr())
	ps, err := pmux.Client(&mux.WsConn{Conn: c}, channel.InitialPMuxConfig(&conf.Cipher))
	if nil != err {
		return nil, err
	}
	return &mux.ProxyMuxSession{Session: ps, NetConn: c}, nil
}

func init() {
	channel.RegisterLocalChannelType("ws", &WebsocketProxy{})
	channel.RegisterLocalChannelType("wss", &WebsocketProxy{})
}
