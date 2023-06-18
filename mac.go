package traefikpluginmanualaccesscontrol

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Key string `json:"key,omitempty"`
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

// entry New function for traefik plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	key, err := NewKey(config.Key)
	if err != nil {
		return nil, err
	}
	return &MAC{
		ctx:    ctx,
		next:   next,
		config: config,
		name:   name,
		key:    key,
	}, nil
}

func (m *MAC) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// set cookie in rw
	rw.Header().Set("Set-Cookie", fmt.Sprintf("tpmac-token=%s", m.key.GenerateToken()))
	m.next.ServeHTTP(rw, req)
}
