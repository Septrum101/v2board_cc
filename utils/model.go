package utils

import (
	"fmt"
	"math/rand"
	"time"
)

type UserAgent struct {
	ID           int
	UA           string
	BannedCounts int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandUA() []UserAgent {
	v := rand.Intn(100)
	UAList := []UserAgent{
		{ID: 0, UA: fmt.Sprintf("ClashX Pro/1.%d.0.2 (com.west2online.ClashXPro; build:1.%d.0.2; macOS 10.13.6) Alamofire/5.4.4", v, v)},
		{ID: 1, UA: fmt.Sprintf("ClashX/1.%d.0 (com.west2online.ClashX; build:1.%d.0; macOS 12.1.0) Alamofire/5.4.4", v, v)},
		{ID: 2, UA: fmt.Sprintf("ClashforWindows/0.%d.%d", rand.Intn(19)+1, rand.Intn(10))},
		{ID: 3, UA: fmt.Sprintf("ClashForAndroid/2.%d.%d.premium", rand.Intn(6), rand.Intn(10))},
		{ID: 4, UA: fmt.Sprintf("Shadowrocket/%d CFNetwork/%d.0.4 Darwin/%d.2.0", rand.Intn(1000)+1000, rand.Intn(1000)+1000, rand.Intn(21))},
		{ID: 5, UA: ""},
	}
	return UAList
}
