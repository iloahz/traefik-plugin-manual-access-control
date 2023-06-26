package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	Server string `json:"server,omitempty" yaml:"server,omitempty" toml:"server,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type TPMAC struct {
	ctx    context.Context
	next   http.Handler
	config *Config
	name   string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	m := &TPMAC{
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

type Request struct {
	IP        string  `json:"ip"`
	UserAgent string  `json:"user_agent"`
	URL       string  `json:"url"`
	Token     *string `json:"token,omitempty"`
}

type Response struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (m *TPMAC) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r := &Request{
		UserAgent: req.UserAgent(),
		URL:       req.Host,
	}
	if tokens := strings.Split(req.RemoteAddr, ":"); len(tokens) == 2 {
		r.IP = tokens[0]
	}
	if ips := req.Header.Values("X-Forwarded-For"); len(ips) > 0 && len(ips[0]) > 0 {
		if tokens := strings.Split(ips[0], ","); len(tokens) > 0 && len(tokens[len(tokens)-1]) > 0 {
			r.IP = tokens[len(tokens)-1]
		}
	}
	if ip := req.Header.Values("X-Real-IP"); len(ip) > 0 && len(ip[0]) > 0 {
		r.IP = ip[0]
	}
	cookie, err := req.Cookie(cookieKey)
	if err != nil || len(cookie.Value) == 0 {
		// generate token
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
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			var t Response
			err = json.Unmarshal(buf, &t)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(rw, &http.Cookie{
				Name:  cookieKey,
				Value: t.Token,
				Path:  "/",
			})
			io.WriteString(rw, fmt.Sprintf("id: %s\nshare this id to admin to get access", t.ID))
			return
		} else {
			// token not generated
			http.Error(rw, "failed to generate token", http.StatusInternalServerError)
			return
		}
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
			// token valid, might refreshed
			defer res.Body.Close()
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			var t Response
			err = json.Unmarshal(buf, &t)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if t.Token != cookie.Value {
				http.SetCookie(rw, &http.Cookie{
					Name:  cookieKey,
					Value: t.Token,
					Path:  "/",
				})
			}
			m.next.ServeHTTP(rw, req)
			return
		} else {
			// token invalid, not allowed, etc.
			http.Error(rw, "not allowed", http.StatusForbidden)
			return
		}
	}
}
