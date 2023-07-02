package clients

import (
	"time"

	"github.com/google/uuid"
	"github.com/iloahz/traefik-plugin-manual-access-control/tokens"
)

type ConsentStatus string

const (
	ConsentStatusDefault ConsentStatus = "default"
	ConsentStatusPending ConsentStatus = "pending"
	ConsentStatusBlocked ConsentStatus = "blocked"
	ConsentStatusAllowed ConsentStatus = "allowed"
)

const (
	AnyIP = "*"
)

type Consent struct {
	ID        string        `json:"id"`
	ClientID  string        `json:"client_id"`
	IP        string        `json:"ip"`
	Host      string        `json:"host"`
	Status    ConsentStatus `json:"status"`
	UpdatedAt int64         `json:"updated_at"`

	token string
}

func NewConsent(clientID string, host string, ip string) *Consent {
	return &Consent{
		ID:        uuid.NewString(),
		ClientID:  clientID,
		Host:      host,
		IP:        ip,
		Status:    ConsentStatusDefault,
		UpdatedAt: time.Now().UnixMilli(),
	}
}

func (c *Consent) Allow() {
	c.Status = ConsentStatusAllowed
	c.UpdatedAt = time.Now().UnixMilli()
}

func (c *Consent) Block() {
	c.Status = ConsentStatusBlocked
	c.UpdatedAt = time.Now().UnixMilli()
}

func (c *Consent) newToken() string {
	c.token = tokens.GenerateToken(c.ClientID, c.Host, c.IP)
	return c.token
}

func (c *Consent) GenerateToken() string {
	if len(c.token) == 0 {
		return c.newToken()
	}
	claims, err := tokens.ValidateToken(c.token)
	if err != nil {
		return c.newToken()
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return c.newToken()
	}
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return c.newToken()
	}
	dur := exp.Sub(iat.Time).Milliseconds()
	now := time.Now().UnixMilli()
	if now > iat.Time.UnixMilli()+dur/2 {
		return c.newToken()
	}
	return c.token
}
