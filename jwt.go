package certmetrics

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

type JWTInspector struct {
	collector Collector
}

func NewJWTInspector(collector Collector) *JWTInspector {
	return &JWTInspector{
		collector: collector,
	}
}

func (i *JWTInspector) InspectToken(token string) error {
	return i.InspectNamedToken(func(claims map[string]interface{}) string {
		sub, ok := claims["sub"]
		if !ok {
			return "unknown"
		}
		return sub.(string)
	}, token)
}

func (i *JWTInspector) InspectNamedToken(subjectName func(claims map[string]interface{}) string, token string) error {
	claims := jwt.MapClaims{}

	parsed, _, err := (&jwt.Parser{}).ParseUnverified(token, &claims)
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	i.collector.StoreCert(i.cert(subjectName(claims), parsed, claims))

	return nil
}

func (i *JWTInspector) cert(subjectName string, _ *jwt.Token, claims jwt.MapClaims) *Cert {
	cert := &Cert{
		Type:    "jwt",
		Subject: subjectName,
	}

	if exp, ok := claims["exp"]; ok {
		cert.ExpiredAt = timeFromUnixString(exp)
	}

	if startedAt, ok := claims["nbf"]; ok {
		cert.StartedAt = timeFromUnixString(startedAt)
	} else if startedAt, ok = claims["iat"]; ok {
		cert.StartedAt = timeFromUnixString(startedAt)
	}

	return cert
}

func timeFromUnixString(value interface{}) time.Time {
	var unix int64 = 0

	switch v := value.(type) {
	case string:
		unix, _ = strconv.ParseInt(v, 10, 64)
	case float64:
		unix = int64(v)
	case int64:
		unix = v
	}

	return time.Unix(unix, 0)
}
