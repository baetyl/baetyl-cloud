package common

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/baetyl/baetyl-go/v2/errors"
)

// Field returns a field
func Field(k string, v interface{}) *F {
	return &F{k, v}
}

// Error returns an error with code and fields
func Error(c Code, fs ...*F) error {
	m := c.String()
	if strings.Contains(m, "{{") {
		vs := map[string]interface{}{}
		for _, f := range fs {
			vs[f.k] = f.v
		}
		t, err := template.New(string(c)).Option("missingkey=zero").Parse(m)
		if err != nil {
			panic(err)
		}
		b := bytes.NewBuffer(nil)
		err = t.Execute(b, vs)
		if err != nil {
			panic(err)
		}
		m = b.String()
	}
	return errors.CodeError(string(c), m)
}

// Field field
type F struct {
	k string
	v interface{}
}
