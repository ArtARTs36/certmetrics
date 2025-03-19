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

import "github.com/artarts36/certmetrics"

func main() {
	certmetrics.JWT.InspectToken("jwt token")
}
```

### Store X509 Certificate

```go
package main

import "github.com/artarts36/certmetrics"

func main() {
	certmetrics.X509.InspectPEMs([]byte("pems"))
}
```
