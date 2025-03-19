package jwtm

import "github.com/artarts36/certmetrics"

type subjectFuncOpt struct {
	namer func(map[string]interface{}) string
}

func WithSubjectNameOf(namer func(claims map[string]interface{}) string) InspectOption {
	return &subjectFuncOpt{namer: namer}
}

func WithoutSubjectName() InspectOption {
	return WithSubjectNameOf(func(map[string]interface{}) string {
		return "unknown"
	})
}

func (o *subjectFuncOpt) apply(claims map[string]interface{}, storing *certmetrics.Cert) {
	storing.Subject = o.namer(claims)
}
