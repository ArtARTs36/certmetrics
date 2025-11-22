package jwtm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/artarts36/certmetrics/pkg/collector"
)

func TestJWTInspectorInspectToken(t *testing.T) {
	tests := []struct {
		Title    string
		Token    string
		Expected []*collector.Cert
	}{
		{
			Title: "John Doe",
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", //nolint:lll
			Expected: []*collector.Cert{
				{
					Type:      "jwt",
					Subject:   "1234567890",
					StartedAt: time.Unix(1516239022, 0),
				},
			},
		},
		{
			Title: "invalid sub type",
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOlsxLDJdLCJuYW1lIjoiSm9obiBEb2UiLCJpYXQiOjE1MTYyMzkwMjJ9.h7TvvcqBC1W6CD8wgf-ZvDxVFSmpsDasnh6o7rOXsQg", //nolint:lll
			Expected: []*collector.Cert{
				{
					Type:      "jwt",
					Subject:   "<invalid>",
					StartedAt: time.Unix(1516239022, 0),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			cl := collector.NewMemoryCollector()
			inspector := NewInspector(cl)

			err := inspector.InspectToken(test.Token)
			require.NoError(t, err)
			assert.Equal(t, test.Expected, cl.Certs())
		})
	}
}
