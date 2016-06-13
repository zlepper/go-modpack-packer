package helpers

import "testing"

type testpair struct {
	s      string
	sep    string
	result bool
}

var tests = []testpair{
	{"string", "rin", true},
	{"STRING", "rin", true},
	{"STRING", "RIN", true},
	{"string", "RIN", true},
	{"other", "rin", false},
}

func TestIgnoreCaseContains(t *testing.T) {
	for _, pair := range tests {
		result := IgnoreCaseContains(pair.s, pair.sep)
		if result != pair.result {
			t.Error("Expected", pair.result, "got", result)
		}
	}
}
