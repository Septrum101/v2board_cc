package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func getUA(ua int) string {
	switch ua {
	case 0:
		v := rand.Intn(100)
		return fmt.Sprintf("ClashX Pro/1.%d.0.2 (com.west2online.ClashXPro; build:1.%d.0.2; macOS 10.13.6) Alamofire/5.4.4", v, v)
	case 1:
		v := rand.Intn(100)
		return fmt.Sprintf("ClashX/1.%d.0 (com.west2online.ClashX; build:1.%d.0; macOS 12.1.0) Alamofire/5.4.4", v, v)
	case 2:
		return fmt.Sprintf("ClashforWindows/0.%d.%d", rand.Intn(19)+1, rand.Intn(10))
	case 3:
		return fmt.Sprintf("ClashForAndroid/2.%d.%d.premium", rand.Intn(6), rand.Intn(10))
	case 4:
		return fmt.Sprintf("Shadowrocket/%d CFNetwork/%d.0.4 Darwin/%d.2.0", rand.Intn(1000)+1000, rand.Intn(1000)+1000, rand.Intn(21))
	default:
		return ""
	}
}
