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
	jwt    *JWT
}

// entry New function for traefik plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	jwt, err := NewJWT(config.Key)
	if err != nil {
		return nil, err
	}
	m := &MAC{
		ctx:    ctx,
		next:   next,
		config: config,
		name:   name,
		jwt:    jwt,
	}
	return m, nil
}

func (m *MAC) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// set cookie in rw
	rw.Header().Set("Set-Cookie", fmt.Sprintf("t-mac-token=%s", m.jwt.GenerateToken()))
	m.next.ServeHTTP(rw, req)
}
