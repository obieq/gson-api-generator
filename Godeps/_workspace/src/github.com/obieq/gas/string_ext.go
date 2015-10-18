package gas

import (
	"bytes"
	"strings"

	"github.com/gedex/inflector"
)

type String string

func (str String) Pluralize() string {
	return inflector.Pluralize(string(str))
}
func (str String) Underscore() string {
	return insertRune(str, '_')
}

func (str String) Dasherize() string {
	return insertRune(str, '-')
}

func insertRune(str String, runeToInsert rune) string {
	buf := bytes.NewBufferString("")
	for i, v := range str {
		if i > 0 && v >= 'A' && v <= 'Z' {
			buf.WriteRune(runeToInsert)
		}
		buf.WriteRune(v)
	}

	return strings.ToLower(buf.String())
}
