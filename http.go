package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

var httpProxy = goproxy.NewProxyHttpServer()

func init() {
	httpProxy.Verbose = true

	httpProxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			newReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
			if err != nil {
				log.Printf("[http] New request error: %v", err)
				return req, nil
			}
			newReq.Header = req.Header

			// 设置自定义拨号器的 HTTP 客户端
			client := &http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return net.Dial(network, addr) // 不强制绑定源IP
					},
				},
			}

			// 发送 HTTP 请求
			resp, err := client.Do(newReq)
			if err != nil {
				log.Printf("[http] Send request error: %v", err)
				return req, nil
			}
			return req, resp
		},
	)
}
