# certmetrics

It's Go library for storing metrics(prometheus) about your X509 certificates or JWT tokens.

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

