package jwtm

import (
	"github.com/artarts36/certmetrics/pkg/collector"
)

type subjectFuncOpt struct {
	namer func(map[string]interface{}) string
}

func WithSubjectNameOf(namer func(claims map[string]interface{}) string) InspectOption {
	return &subjectFuncOpt{namer: namer}
}

func WithoutSubjectName() InspectOption {
	return WithSubjectNameOf(func(map[string]interface{}) string {
		return "<hidden>"
	})
}

func (o *subjectFuncOpt) apply(claims map[string]interface{}, storing *collector.Cert) {
	storing.Subject = o.namer(claims)
}
