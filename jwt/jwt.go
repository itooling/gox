package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strings"
	"time"

	"github.com/itooling/gox/sys"
)

const (
	ALG = "HS256"
	JWT = "JWT"
)

var (
	secret  string
	expires int
)

func init() {
	secret = sys.String("token.secret")
	if secret == "" {
		secret = "secret"
	}
	expires = sys.Int("token.expires")
	if expires == 0 {
		expires = 30
	}
}

type Header struct {
	Alg string `json:"alg,omitempty"`
	Typ string `json:"typ,omitempty"`
}

type Payload struct {
	JwtId     string         `json:"jti,omitempty"`
	Issuer    string         `json:"iss,omitempty"`
	Subject   string         `json:"sub,omitempty"`
	IssuedAt  int64          `json:"iat,omitempty"`
	Audience  string         `json:"aud,omitempty"`
	ExpiresAt int64          `json:"exp,omitempty"`
	NotBefore int64          `json:"nbf,omitempty"`
	Data      map[string]any `json:"dat,omitempty"`
}

type Signature struct {
	header          Header
	payload         Payload
	headerString    string
	payloadString   string
	signatureString string
}

func encode(o any) string {
	js, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(js)
}

func (s *Signature) Sign(secret string) {
	s.headerString = encode(s.header)
	s.payloadString = encode(s.payload)

	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(secret))

	h.Write([]byte(s.headerString + "." + s.payloadString))
	res := h.Sum(nil)

	s.signatureString = base64.URLEncoding.EncodeToString(res)
}

func (s *Signature) String() string {
	return fmt.Sprintf("%s.%s.%s", s.headerString, s.payloadString, s.signatureString)
}

func Create(id string) string {
	s := Signature{
		header: Header{
			Alg: ALG,
			Typ: JWT,
		},
		payload: Payload{
			JwtId:     id,
			ExpiresAt: time.Now().Add(time.Duration(expires) * time.Minute).Unix(),
		},
	}
	s.Sign(secret)
	return s.String()
}

func Creates(id string, data map[string]any) string {
	s := Signature{
		header: Header{
			Alg: ALG,
			Typ: JWT,
		},
		payload: Payload{
			JwtId:     id,
			Data:      data,
			ExpiresAt: time.Now().Add(time.Duration(expires) * time.Minute).Unix(),
		},
	}
	s.Sign(secret)
	return s.String()
}

func Verify(token string) bool {
	sp := strings.Split(token, ".")
	if len(sp) == 3 {
		h := sp[0]
		p := sp[1]
		s := sp[2]

		pp, err := base64.URLEncoding.DecodeString(p)
		if err != nil {
			return false
		}
		pd := Payload{}
		err = json.Unmarshal(pp, &pd)
		if err != nil {
			return false
		}
		expires := time.Now().Before(time.Unix(pd.ExpiresAt, 0))

		if expires {
			sh := hmac.New(
				func() hash.Hash { return sha256.New() },
				[]byte(secret),
			)

			sh.Write([]byte(h + "." + p))
			ss := sh.Sum(nil)

			return s == base64.URLEncoding.EncodeToString(ss)
		}
	}
	return false
}

func GetPayload(token string) *Payload {
	sp := strings.Split(token, ".")
	if len(sp) == 3 {
		p := sp[1]
		pp, err := base64.URLEncoding.DecodeString(p)
		if err != nil {
			return nil
		}
		pd := Payload{}
		err = json.Unmarshal(pp, &pd)
		if err != nil {
			return nil
		}
		return &pd
	}
	return nil
}
