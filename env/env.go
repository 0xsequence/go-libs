// Package env defines our application environments (local, test, dev, staging, prod)
// with utilities for comparison, string conversion, and text marshaling.
package env

import (
	"fmt"
	"slices"
	"strings"
)

type Env uint8

const (
	EnvLocal Env = iota
	EnvTest
	EnvDev
	EnvDev2
	EnvNext
	EnvProd
)

var environments = []string{
	"local", // 0
	"test",  // 1
	"dev",   // 2
	"dev2",  // 3
	"next",  // 4
	"prod",  // 5
}

func (e Env) Is(envs ...Env) bool {
	return slices.Contains(envs, e)
}

func (e Env) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Env) UnmarshalText(text []byte) error {
	enum := string(text)

	// on empty string fallback to "local"
	if enum == "" {
		*e = EnvLocal
		return nil
	}

	for i, name := range environments {
		if enum == name {
			*e = Env(i)
			return nil
		}
	}

	return fmt.Errorf("unknown env=(%s), supported=(%s)", text, strings.Join(environments, ","))
}

func (e Env) String() string {
	if int(e) >= len(environments) {
		return fmt.Sprintf("Env(%d)", e)
	}

	return environments[e]
}
