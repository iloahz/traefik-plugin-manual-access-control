package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/api"
	"github.com/iloahz/traefik-plugin-manual-access-control/clients"
	"github.com/iloahz/traefik-plugin-manual-access-control/logs"
)

func main() {
	logs.Init("debug")
	if os.Getenv("DEBUG") == "true" {
		go func() {
			time.Sleep(time.Second)
			clients.GetClient("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36", "116.179.32.218", "cdn.home.iloahz.com")

			time.Sleep(time.Second)
			c2 := clients.GetClient("Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1", "221.192.199.49", "code.home.iloahz.com")
			c2.Allow("code.home.iloahz.com", clients.AnyIP)

			time.Sleep(time.Second)
			c3 := clients.GetClient("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36", "180.163.220.66", "chatgpt.home.iloahz.com")
			c3.Block("chatgpt.home.iloahz.com", clients.AnyIP)
		}()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	api.CreateServer()
}
