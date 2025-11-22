package collector

type MemoryCollector struct {
	certs []*Cert
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{
		certs: []*Cert{},
	}
}

func (c *MemoryCollector) StoreCert(cert *Cert) {
	c.certs = append(c.certs, cert)
}

func (c *MemoryCollector) Certs() []*Cert {
	return c.certs
}
