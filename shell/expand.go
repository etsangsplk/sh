// Copyright (c) 2018, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package shell

import (
	"context"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Expand performs shell expansion on s, using env to resolve variables.
// The expansion will apply to parameter expansions like $var and
// ${#var}, but also to arithmetic expansions like $((var + 3)), and
// command substitutions like $(echo foo).
//
// If env is nil, the current environment variables are used.
//
// Any side effects or modifications to the system are forbidden when
// interpreting the program. This is enforced via whitelists when
// executing programs and opening paths. The interpreter also has a timeout of
// two seconds.
func Expand(s string, env func(string) string) (string, error) {
	p := syntax.NewParser()
	word, err := p.Document(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	r := pureRunner()
	if env != nil {
		r.Env = interp.FuncEnviron(env)
	}
	ctx, cancel := context.WithTimeout(context.Background(), pureRunnerTimeout)
	defer cancel()
	fields, err := r.Fields(ctx, word)
	return strings.Join(fields, ""), err
}
