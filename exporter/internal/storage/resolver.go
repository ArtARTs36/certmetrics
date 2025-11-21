package storage

import (
	"regexp"
)

type Resolver struct {
	rules []*ResolveRule
	def   Storage
}

type ResolveRule struct {
	Regex   *regexp.Regexp
	Storage Storage
}

func NewResolver(
	def Storage,
	rules []*ResolveRule,
) *Resolver {
	return &Resolver{
		rules: rules,
		def:   def,
	}
}

func (r *Resolver) Resolve(path string) Storage {
	for _, rule := range r.rules {
		if rule.Regex.MatchString(path) {
			return rule.Storage
		}
	}

	return r.def
}
