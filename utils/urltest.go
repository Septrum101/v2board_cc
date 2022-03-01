package utils

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/go-resty/resty/v2"
)

func URLTest(p *Nodes) (aliveProxies Nodes, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseURL, err := url.Parse("http://www.gstatic.com/generate_204")
	if err != nil {
		return
	}

	addr, err := urlToMetadata(baseURL.String())
	if err != nil {
		return
	}

	instance, err := p.DialContext(ctx, &addr)
	if err != nil {
		return
	}
	defer func(instance C.Conn) {
		err = instance.Close()
		if err != nil {
			return
		}
	}(instance)

	transport := &http.Transport{DialContext: func(context.Context, string, string) (net.Conn, error) { return instance, nil }}

	resp, err := resty.New().
		SetTransport(transport).
		R().SetContext(ctx).Head(baseURL.String())

	if resp.StatusCode() == 204 {
		aliveProxies = *p
	}
	return
}
