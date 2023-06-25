package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Server string `json:"server,omitempty" yaml:"server,omitempty" toml:"server,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type MAC struct {
	ctx    context.Context
	next   http.Handler
	config *Config
	name   string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	m := &MAC{
		ctx:    ctx,
		next:   next,
		config: config,
		name:   name,
	}
	return m, nil
}

const (
	cookieKey = "tpmac"
)

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrNotAllowed   = fmt.Errorf("not allowed")
)

type Request struct {
	IP        string  `json:"ip"`
	UserAgent string  `json:"user_agent"`
	URL       string  `json:"url"`
	Token     *string `json:"token,omitempty"`
}

func (m *MAC) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieKey)
	r := &Request{
		IP:        req.RemoteAddr,
		UserAgent: req.UserAgent(),
		URL:       req.URL.String(),
	}
	if err != nil {
		// generate token
		r.Token = &cookie.Value
		buf, err := json.Marshal(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := http.Post(m.config.Server+"/api/token/generate", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.StatusCode == http.StatusOK {
			// token generated
			defer res.Body.Close()
			type Response struct {
				Token string `json:"token"`
			}
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			var token Response
			err = json.Unmarshal(buf, &token)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/", cookieKey, token.Token))
			m.next.ServeHTTP(rw, req)
		} else {
			// token not generated
			http.Error(rw, ErrNotAllowed.Error(), http.StatusForbidden)
		}
		return
	} else {
		// validate token
		r.Token = &cookie.Value
		buf, err := json.Marshal(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := http.Post(m.config.Server+"/api/token/validate", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.StatusCode == http.StatusOK {
			// token valid
			m.next.ServeHTTP(rw, req)
			return
		} else {
			// token invalid
			http.Error(rw, ErrInvalidToken.Error(), http.StatusForbidden)
			return
		}
	}
}
