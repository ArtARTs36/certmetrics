package certmetrics

type Collector interface {
	StoreCert(cert *Cert)
}
