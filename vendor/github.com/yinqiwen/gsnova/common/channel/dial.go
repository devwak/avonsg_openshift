package channel

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/yinqiwen/gsnova/common/dns"
	"github.com/yinqiwen/gsnova/common/helper"
	"github.com/yinqiwen/gsnova/common/hosts"
	"github.com/yinqiwen/gsnova/common/logger"
	"github.com/yinqiwen/gsnova/common/netx"
	"github.com/yinqiwen/gsnova/common/protector"
)

func NewTLSConfig(conf *ProxyChannelConfig) *tls.Config {
	tlscfg := &tls.Config{}
	logger.Info("SecureVerify : %v", conf.SecureVerify)
	tlscfg.InsecureSkipVerify = !conf.SecureVerify //true
	if len(conf.SNI) > 0 {
		tlscfg.ServerName = conf.SNI[0]
	}
	return tlscfg
}

func DialServerByConf(server string, conf *ProxyChannelConfig) (net.Conn, error) {
	rurl, err := url.Parse(server)
	if nil != err {
		return nil, err
	}
	hostport := rurl.Host
	tcpHost, tcpPort, err := net.SplitHostPort(hostport)
	if nil != err {
		switch rurl.Scheme {
		case "http", "ws", "tcp", "tcp4", "tcp6":
			tcpHost = rurl.Host
			tcpPort = "80"
		case "ssh":
			tcpPort = "22"
			tcpHost = rurl.Host
		case "http2", "https", "quic", "kcp", "tls", "wss":
			tcpHost = rurl.Host
			tcpPort = "443"
		default:
			return nil, fmt.Errorf("Invalid scheme:%s", rurl.Scheme)
		}
		hostport = net.JoinHostPort(tcpHost, tcpPort)
	}
	tlscfg := NewTLSConfig(conf)
	if len(tlscfg.ServerName) == 0 {
		if net.ParseIP(tcpHost) == nil {
			tlscfg.ServerName = tcpHost
		}
	}

	if len(conf.SNIProxy) > 0 && tcpPort == "443" {
		if net.ParseIP(conf.SNIProxy) == nil {
			if hosts.InHosts(conf.SNIProxy) {
				hostport = hosts.GetAddr(conf.SNIProxy, "443")
				tcpHost, _, _ = net.SplitHostPort(hostport)
			} else {
				logger.Info("SNIProxy Not exist in hosts:%s", conf.SNIProxy)
			}
		} else {
			tcpHost = conf.SNIProxy
			hostport = net.JoinHostPort(tcpHost, tcpPort)
		}
		logger.Info("Try to connect %s via sni proxy:%s", server, hostport)
	}

	var conn net.Conn
	dailTimeout := conf.LocalDialMSTimeout
	if 0 == dailTimeout {
		dailTimeout = 5000
	}
	timeout := time.Duration(dailTimeout) * time.Millisecond
	connAddr := hostport
	if len(conf.Proxy) == 0 {
		if net.ParseIP(tcpHost) == nil {
			iphost, err := dns.DnsGetDoaminIP(tcpHost)
			if nil != err {
				return nil, err
			}
			hostport = net.JoinHostPort(iphost, tcpPort)
		}
		if len(conf.P2PToken) > 0 {
			opt := &protector.NetOptions{
				ReusePort:   protector.SupportReusePort(),
				DialTimeout: timeout,
			}
			conn, err = protector.DialContextOptions(context.Background(), "tcp", hostport, opt)
		} else {
			conn, err = netx.DialTimeout("tcp", hostport, timeout)
		}

	} else {
		if len(conf.P2PToken) > 0 {
			conn, err = helper.ProxyDial(conf.Proxy, "", hostport, timeout, true)
		} else {
			conn, err = helper.ProxyDial(conf.Proxy, "", hostport, timeout, false)
		}
		connAddr = conf.Proxy
	}
	if nil == err {
		switch rurl.Scheme {
		case "tls":
			fallthrough
		case "http2":
			if conf.ForceTls13 == "tls13" {
				tlscfg.MaxVersion = tls.VersionTLS13
				tlscfg.MinVersion = tls.VersionTLS13
				tlscfg.CipherSuites = []uint16{
					tls.TLS_AES_128_GCM_SHA256,
					tls.TLS_AES_256_GCM_SHA384,
					tls.TLS_CHACHA20_POLY1305_SHA256,
				}
				logger.Info("Http2 or tls ForceTls13, Now tls13 dial.")
			} else if conf.ForceTls13 == "tls12" {
				tlscfg.MaxVersion = tls.VersionTLS12
				tlscfg.MinVersion = tls.VersionTLS12
				tlscfg.CipherSuites = []uint16{
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
				logger.Info("Http2 or tls ,tls12 dial.")
			} else { // "" or "auto" or other
				tlscfg.MaxVersion = tls.VersionTLS13
				tlscfg.MinVersion = tls.VersionTLS12
				tlscfg.CipherSuites = []uint16{
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
				logger.Info("Http2 or tls ,tls12-tls13 auto dial.")
			}
			tlsconn := tls.Client(conn, tlscfg)
			err = tlsconn.Handshake()
			if err != nil {
				logger.Notice("TLS Handshake Failed %v", err)
				return nil, err
			}
			conn = tlsconn
		}
	}
	if nil != err {
		logger.Notice("Connect %s failed with reason:%v.", server, err)
	} else {
		logger.Debug("Connect %s success via %s.", server, connAddr)
	}
	return conn, err
}

func NewDialByConf(conf *ProxyChannelConfig, scheme string) func(network, addr string) (net.Conn, error) {
	localDial := func(network, addr string) (net.Conn, error) {
		//log.Printf("Connect %s", addr)
		server := fmt.Sprintf("%s://%s", scheme, addr)
		return DialServerByConf(server, conf)
	}
	return localDial
}

var httpClientMap sync.Map

func NewHTTPClient(conf *ProxyChannelConfig, scheme string) (*http.Client, error) {
	tr := &http.Transport{
		Dial:                  NewDialByConf(conf, scheme),
		DisableCompression:    true,
		MaxIdleConnsPerHost:   2 * int(conf.ConnsPerServer),
		ResponseHeaderTimeout: time.Duration(conf.HTTP.ReadTimeout) * time.Millisecond,
	}
	// if len(conf.SNI) > 0 {
	// 	tlscfg := &tls.Config{}
	// 	tlscfg.InsecureSkipVerify = true
	// 	tlscfg.ServerName = conf.SNI[0]
	// 	tr.TLSClientConfig = tlscfg
	// }
	// if len(conf.Proxy) > 0 {
	// 	proxyUrl, err := url.Parse(conf.Proxy)
	// 	if nil != err {
	// 		logger.Error("[ERROR]Invalid proxy url:%s to create http client.", conf.Proxy)
	// 		return nil, err
	// 	}
	// 	tr.Proxy = http.ProxyURL(proxyUrl)
	// }
	hc := &http.Client{}
	//hc.Timeout = tr.ResponseHeaderTimeout
	hc.Transport = tr
	localClient, loaded := httpClientMap.LoadOrStore(conf, hc)
	if loaded {
		return localClient.(*http.Client), nil
	}
	return hc, nil
}
