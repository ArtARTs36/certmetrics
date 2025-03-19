package certmetrics

import "time"

type Cert struct {
	ID string

	Type    string
	Subject string

	StartedAt time.Time
	ExpiredAt time.Time
}
