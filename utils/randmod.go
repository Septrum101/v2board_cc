package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func randUA() string {
	userAgent := []string{"axxray", "clash", "passwall", "quantumult%20x", "ssrplus", "shadowrocket", "shadowsocks", "stash", "surfboard", "surge", "v2rayn", "v2rayng"}
	return userAgent[rand.Intn(len(userAgent))]
}
