package services

import "testing"

type testpair struct {
	input  string
	output string
}

var tests = []testpair{
	{"Влад бумага", "vlad-bumaga"},
	{"Николай соболев", "nikolay-sobolev"},
	{"Гусейн Гасанов", "guseyn-gasanov"},
	{"Прохор Шаляпин", "prohor-shalyapin"},
}

func TestChangeWord(t *testing.T) {
	for _, pair := range tests {
		v := changeWord(pair.input)
		if v != pair.output {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}
