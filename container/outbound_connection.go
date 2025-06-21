package container

import (
	"context"
	"github.com/web-page-analysis/bootstrap"
	"net"
	"net/http"
	"time"
)

var (
	connectionClient ConnectionClientConfig
)

type ConnectionClientConfig struct {
	HttpClientDefault http.Client
}

type outBoundConnection struct {
	outboundConf bootstrap.OutboundConfig
}

func (o outBoundConnection) Get(ctx context.Context, url string) (*http.Response, error) {
	client := connectionClient.HttpClientDefault
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

type OutBoundConnection interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

func InitOutBoundConnection(conf bootstrap.Config) OutBoundConnection {

	initConnection(conf.OutboundConf)
	return &outBoundConnection{
		outboundConf: conf.OutboundConf,
	}
}

func initConnection(timeoutConf bootstrap.OutboundConfig) {

	connectionClient.HttpClientDefault = getHttpClient(timeoutConf)
}

func getHttpClient(to bootstrap.OutboundConfig) http.Client {
	return http.Client{
		Timeout: time.Millisecond * time.Duration(to.DialTimeout),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Millisecond * time.Duration(to.DialTimeout),
			}).DialContext,
		},
	}
}
