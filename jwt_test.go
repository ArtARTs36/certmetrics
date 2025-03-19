package certmetrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTInspectorInspectToken(t *testing.T) {
	tests := []struct {
		Title    string
		Token    string
		Expected []*Cert
	}{
		{
			Title: "John Doe",
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			Expected: []*Cert{
				{
					Type:      "jwt",
					Subject:   "1234567890",
					StartedAt: time.Unix(1516239022, 0),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			collector := NewMemoryCollector()
			inspector := NewJWTInspector(collector)

			err := inspector.InspectToken(test.Token)
			require.NoError(t, err)
			assert.Equal(t, test.Expected, collector.certs)
		})
	}
}
