package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func randUA() string {
	var ua []string
	//ua = append(ua, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36 Edg/97.0.1072.62")
	//ua = append(ua, "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0")
	//ua = append(ua, "Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7")
	//ua = append(ua, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36 OPR/77.0.4054.277")
	//ua = append(ua, "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0")
	//ua = append(ua, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	//ua = append(ua, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	ua = append(ua , "axxray")
	ua = append(ua , "clash")
	ua = append(ua , "passwall")
	ua = append(ua , "quantumult%20x")
	ua = append(ua , "ssrplus")
	ua = append(ua , "shadowrocket")
	ua = append(ua , "shadowsocks")
	ua = append(ua , "stash")
	ua = append(ua , "surfboard")
	ua = append(ua , "surge")
	ua = append(ua , "v2rayn")
	ua = append(ua , "v2rayng")
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(ua) - 1)
	return ua[idx]
}

func randStr(length int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func randEmail() string {
	mailProvider := []string{"qq.com", "gmail.com", "sina.com", "126.com", "sina.cn"}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	idx := rand.Intn(len(mailProvider) - 1)
	mailHost := mailProvider[idx]
	return fmt.Sprintf("%s@%s", randStr(10), mailHost)
}
