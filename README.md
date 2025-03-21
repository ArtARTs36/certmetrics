# certmetrics

It's Go library and exporter for storing metrics(prometheus) about your X509 certificates or JWT tokens.

```
# HELP certmetrics_cert_expiry_days Days until certificate expiry
# TYPE certmetrics_cert_expiry_days gauge
certmetrics_cert_expiry_days{id="my-name",subject="Russian Trusted Sub CA"} 1581
certmetrics_cert_expiry_days{id="super-name",subject="Russian Trusted Sub CA"} 715
# HELP certmetrics_cert_info Cert info
# TYPE certmetrics_cert_info gauge
certmetrics_cert_info{expired_at="2027-03-06 11:25:19",id="my-name",reporter="exporter",started_at="2022-03-02 11:25:19",subject="Russian Trusted Sub CA",type="x509"} 1
certmetrics_cert_info{expired_at="2027-03-06 11:25:19",id="super-name",reporter="exporter",started_at="2022-03-02 11:25:19",subject="Russian Trusted Sub CA",type="x509"} 1
certmetrics_cert_info{expired_at="<unknown>",id="gpt-token",reporter="exporter",started_at="2018-01-18 04:30:22",subject="1234567890",type="jwt"} 1
```

## Usage

### 🐳 As exporter

Use docker image `artarts36/certmetrics:v0.1.0`

1. Copy [config](./exporter/certmetrics.yaml) as `certmetrics.yaml`
2. Run `docker run -v ./certmetrics.yaml:/app/certmetrics.yaml -p 8010:8010 artarts36/certmetrics:v0.1.0`

### 📚 As Go library

Run: `go get github.com/artarts36/certmetrics`

1. Register default collectors
2. Store certs and tokens.

```go
package main

import (
	"github.com/artarts36/certmetrics"
	"github.com/artarts36/certmetrics/jwtm"
	"github.com/artarts36/certmetrics/x509m"
)

func main() {
	certmetrics.Register()
	
	jwtm.InspectToken("token", jwtm.WithID("super-token"))
	x509m.InspectPEMs([]byte("pems"), x509m.WithID("super-ca"))
}
```
