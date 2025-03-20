# certmetrics

It's Go library for storing metrics(prometheus) about your X509 certificates or JWT tokens.

```
go_cert_info{expired_at="2027-03-06 11:25:19",id="russian-ca",started_at="2022-03-02 11:25:19",subject="Russian Trusted Sub CA",type="x509"} 1
go_cert_info{expired_at="2032-02-27 21:04:15",id="russian-ca",started_at="2022-03-01 21:04:15",subject="Russian Trusted Root CA",type="x509"} 1
go_cert_info{expired_at="<unknown>",id="gpt-token",started_at="2018-01-18 04:30:22",subject="1234567890",type="jwt"} 1
```

Register default collectors with easy:

```go
package main

import "github.com/artarts36/certmetrics"

func main() {
	certmetrics.Register()
}
```

## Usage examples

### Store JWT Token

```go
package main

import "github.com/artarts36/certmetrics/jwtm"

func main() {
	jwtm.InspectToken("jwt token")
}
```

### Store X509 Certificate

```go
package main

import "github.com/artarts36/certmetrics/x509m"

func main() {
	x509m.InspectPEMs([]byte("pems"))
}
```

### Store with ID label 

JWT

```go
package main

import (
	"github.com/artarts36/certmetrics/jwtm"
)

func main() {
	jwtm.InspectToken("token", jwtm.WithID("super-token"))
}

```
X509

```go
package main

import "github.com/artarts36/certmetrics/x509m"

func main() {
	x509m.InspectPEMs([]byte("pems"), x509m.WithID("super-ca"))
}
```

