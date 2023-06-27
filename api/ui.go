package api

import "github.com/gin-contrib/static"

func ServeUI() {
	r.Use(static.Serve("/", static.LocalFile("ui/dist", false)))
}
