package main

import (
	"testing"
    "math/rand"
)


func randBinaryString(n int) string {
    var binRunes = []rune("01")
    out := make([]rune, n)
    for i := range out {
        out[i] = binRunes[rand.Intn(len(binRunes))]
    }
    return string(out)
}


func TestNewDomainRecord(t *testing.T) {
	var criteria = []struct {
		Name     string
		Domain   string
		Negative bool
	}{
		{"Given a valid domain", "godaddy.com", false},
		{"Given a domain without a TLD", "localhost", false},
		{"Given an empty domain", "", true},
		{"Given a long name", randBinaryString(8193), true},
	}
	for _, test := range criteria {
		t.Run(test.Name, func(t *testing.T) {
			if _, err := NewDomainRecord(test.Domain, "A", "127.0.0.1", 60, 0); err != nil {
				if !test.Negative {
					t.Errorf("failed to create new domain record: %s", err)
				}
			}
		})
	}
}
