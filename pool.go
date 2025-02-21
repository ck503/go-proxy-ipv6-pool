package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"sync"
)

type IPPool struct {
	pool *sync.Pool
	CIDR string // cidr
	//Eth              string // 网卡名称
	IPV6Net          *net.IPNet
	MaskSize         int // 获取网络部分和掩码长度
	RandomPartLength int // 随机段长度
}

func NewIPPool(initSize int, cidr string) *IPPool {
	p := &IPPool{
		CIDR: cidr,
		//Eth:  ethName,
	}

	_, ipv6Net, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("CIDR error:%s\n", err)
		os.Exit(1)
	}

	p.IPV6Net = ipv6Net
	p.MaskSize, _ = ipv6Net.Mask.Size()
	p.RandomPartLength = 128 - p.MaskSize

	p.pool = &sync.Pool{
		New: func() interface{} {
			ip6, _ := p.generateRandomIPv6()
			return ip6
		},
	}

	// init size
	for i := 0; i < initSize; i++ {
		ip6, _ := p.generateRandomIPv6()
		p.pool.Put(ip6)
	}

	return p
}

func (p *IPPool) Get() string {
	return p.pool.Get().(string)
}

func (p *IPPool) Put(s string) {
	p.pool.Put(s)
}

// gen ipv6 and register linux proxy_ndp
func (p *IPPool) generateRandomIPv6() (string, error) {
	// 生成随机部分
	randomPart := make([]byte, p.RandomPartLength/8)
	_, err := rand.Read(randomPart)
	if err != nil {
		return "", err
	}

	// 获取网络部分
	networkPart := p.IPV6Net.IP.To16()

	// 合并网络部分和随机部分
	for i := 0; i < len(randomPart); i++ {
		networkPart[16-len(randomPart)+i] = randomPart[i]
	}

	var res = networkPart.String()

	// add to linux core proxy_ndp
	//p.addProxyIP(res)

	return res, nil
}

//func (p *IPPool) addProxyIP(ip string) {
//	cmd := exec.Command("ip", "neigh", "add", "proxy", ip, "dev", p.Eth)
//	err := cmd.Run()
//	if err != nil {
//		fmt.Println("add proxy ip to linux core err", err)
//	}
//}
