package collector

type Collector interface {
	StoreCert(cert *Cert)
}
