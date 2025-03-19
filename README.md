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

## Store JWT Token

```go
package main

import "github.com/artarts36/certmetrics"

func main() {
	certmetrics.JWT.InspectToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
}
```
