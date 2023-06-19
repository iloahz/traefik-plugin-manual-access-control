package traefik_plugin_manual_access_control

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type MAC struct {
	ctx    context.Context
	next   http.Handler
	config *Config
	name   string
	key    *Key
}

// func (m *MAC) createServer() {
// 	r := gin.Default()
// 	r.Static("/", "./ui/dist")
// 	server := &http.Server{
// 		Addr:    fmt.Sprintf(":%d", m.config.Port),
// 		Handler: r,
// 	}
// 	go func() {
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			fmt.Printf("listen: %s\n", err)
// 		}
// 	}()
// 	go func() {
// 		<-m.ctx.Done()
// 		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
// 		defer shutdownCancel()
// 		server.Shutdown(shutdownCtx)
// 	}()
// }

// entry New function for traefik plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	key, err := NewKey(config.Key)
	if err != nil {
		return nil, err
	}
	m := &MAC{
		ctx:    ctx,
		next:   next,
		config: config,
		name:   name,
		key:    key,
	}
	// m.createServer()
	return m, nil
}

func (m *MAC) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// set cookie in rw
	rw.Header().Set("Set-Cookie", fmt.Sprintf("tpmac-token=%s", m.key.GenerateToken()))
	m.next.ServeHTTP(rw, req)
}
