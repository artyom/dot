// Package dot provides some known DNS-over-TLS (DOT) resolvers.
package dot

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"time"
)

// Cloudflare returns Resolver that uses Cloudflare service on 1.1.1.1 and
// 1.0.0.1 on port 853.
//
// See https://developers.cloudflare.com/1.1.1.1/dns-over-tls/ for details.
func Cloudflare() *net.Resolver {
	return newResolver("cloudflare-dns.com", "1.1.1.1:853", "1.0.0.1:853")
}

// Quad9 returns Resolver that uses Quad9 service on 9.9.9.9 and 149.112.112.112
// on port 853.
//
// See https://quad9.net/faq/ for details.
func Quad9() *net.Resolver {
	return newResolver("dns.quad9.net", "9.9.9.9:853", "149.112.112.112:853")
}

// Google returns Resolver that uses Google Public DNS service on 8.8.8.8 and
// 8.8.4.4 on port 853.
//
// See https://developers.google.com/speed/public-dns/ for details.
func Google() *net.Resolver {
	return newResolver("dns.google", "8.8.8.8:853", "8.8.4.4:853")
}

// LibreOps returns Resolver that uses LibreDNS service on 116.203.115.192 on
// port 853 operated by LibreOps.
//
// See https://libredns.gr/ for details.
func LibreOps() *net.Resolver {
	return newResolver("dot.libredns.gr", "116.203.115.192:853")
}

func newResolver(serverName string, addrs ...string) *net.Resolver {
	if serverName == "" {
		panic("dot: server name cannot be empty")
	}
	if len(addrs) == 0 {
		panic("dot: addrs cannot be empty")
	}
	var d net.Dialer
	cfg := &tls.Config{
		ServerName:         serverName,
		ClientSessionCache: tls.NewLRUClientSessionCache(0),
	}
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			conn, err := d.DialContext(ctx, "tcp", addrs[rand.Intn(len(addrs))])
			if err != nil {
				return nil, err
			}
			conn.(*net.TCPConn).SetKeepAlive(true)
			conn.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute)
			return tls.Client(conn, cfg), nil
		},
	}
}
