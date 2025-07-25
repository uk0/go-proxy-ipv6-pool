package main

import (
	"context"
	"log"
	"net"

	socks5 "github.com/armon/go-socks5"
)

var socks5Conf = &socks5.Config{}
var socks5Server *socks5.Server

func SocketInit(socksUser string, socksPass string) {
	// 指定出口 IP 地址
	// 指定本地出口 IPv6 地址

	creds := socks5.StaticCredentials{
		socksUser: socksPass,
	}
	auth := socks5.UserPassAuthenticator{Credentials: creds}
	// 创建一个 SOCKS5 服务器配置
	socks5Conf = &socks5.Config{
		AuthMethods: []socks5.Authenticator{auth},
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {

			outgoingIP, err := generateRandomIPv6(cidr)
			if err != nil {
				log.Printf("Generate random IPv6 error: %v", err)
				return nil, err
			}
			outgoingIP = "[" + outgoingIP + "]"

			// 使用指定的出口 IP 地址创建连接
			localAddr, err := net.ResolveTCPAddr("tcp", outgoingIP+":0")
			if err != nil {
				log.Printf("[socks5] Resolve local address error: %v", err)
				return nil, err
			}
			dialer := net.Dialer{
				LocalAddr: localAddr,
			}
			// 通过代理服务器建立到目标服务器的连接

			log.Println("[socks5]", addr, "via", outgoingIP)
			return dialer.DialContext(ctx, network, addr)
		},
	}
	var err error
	// 创建 SOCKS5 服务器
	socks5Server, err = socks5.New(socks5Conf)
	if err != nil {
		log.Fatal(err)
	}
}
