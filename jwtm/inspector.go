package jwtm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/artarts36/certmetrics"
)

type Inspector struct {
	collector certmetrics.Collector
}

func NewInspector(collector certmetrics.Collector) *Inspector {
	return &Inspector{
		collector: collector,
	}
}

func (i *Inspector) InspectToken(token string, opts ...InspectOption) error {
	claims := jwt.MapClaims{}

	_, _, err := (&jwt.Parser{}).ParseUnverified(token, &claims)
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	i.collector.StoreCert(i.cert(claims, opts))

	return nil
}

func (i *Inspector) InspectClaims(claims map[string]interface{}, opts ...InspectOption) {
	i.collector.StoreCert(i.cert(claims, opts))
}

func (i *Inspector) cert(claims jwt.MapClaims, opts []InspectOption) *certmetrics.Cert {
	cert := &certmetrics.Cert{
		Type: "jwt",
	}

	if subjectName, ok := claims["sub"]; ok {
		if sn, snok := subjectName.(string); snok {
			cert.Subject = sn
		} else {
			cert.Subject = "<invalid>"
		}
	} else {
		cert.Subject = "unknown"
	}

	if exp, ok := claims["exp"]; ok {
		cert.ExpiredAt = timeFromUnixString(exp)
	}

	if startedAt, ok := claims["nbf"]; ok {
		cert.StartedAt = timeFromUnixString(startedAt)
	} else if startedAt, ok = claims["iat"]; ok {
		cert.StartedAt = timeFromUnixString(startedAt)
	}

	for _, opt := range opts {
		opt.apply(claims, cert)
	}

	return cert
}

func timeFromUnixString(value interface{}) time.Time {
	var unix int64

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
