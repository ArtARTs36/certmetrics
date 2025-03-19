package certmetrics

import "time"

type Cert struct {
	Type    string
	Subject string

	StartedAt time.Time
	ExpiredAt time.Time
}
