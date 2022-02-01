package config

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies"`
}

type CFG struct {
	V2boardDomain string `yaml:"v2BoardDomain"`
	OriginIP      string `yaml:"originIP"`
	Connections   int    `yaml:"connections"`

}

var Cfg CFG

func init() {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Unmarshal(f, &Cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{}
	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}
	return rawCfg, nil
}

func ParseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, err error) {
	proxies = make(map[string]C.Proxy)
	proxiesConfig := cfg.Proxy
	for _, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			continue
		}
		if _, exist := proxies[proxy.Name()]; exist {
			u4, _ := uuid.NewV4()
			proxies[fmt.Sprintf("%s[%s]", proxy.Name(), strings.ReplaceAll(u4.String(), "-", ""))] = proxy
			continue
		}
		proxies[proxy.Name()] = proxy
	}
	fmt.Printf("total nodes: %d\n", len(proxies))
	return
}
