package main

import (
	"testing"
)

func TestNewDomainRecord(t *testing.T) {
	var criteria = []struct {
		Name     string
		Domain   string
		Negative bool
	}{
		{"Given a valid domain", "godaddy.com", false},
		{"Given a domain without a TLD", "localhost", false},
		{"Given an empty domain", "", true},
		{"Given a long name", "Lopado­temacho­selacho­galeo­kranio­leipsano­drim­hypo­trimmato­silphio­parao­melito­katakechy­meno­kichl­epi­kossypho­phatto­perister­alektryon­opte­kephallio­kigklo­peleio­lagoio­siraio­baphe­tragano­pterygon", true},
	}
	for _, test := range criteria {
		t.Run(test.Name, func(t *testing.T) {
			if _, err := NewDomainRecord(test.Domain, "A", "127.0.0.1", 60); err != nil {
				if !test.Negative {
					t.Errorf("failed to create new domain record: %s", err)
				}
			}
		})
	}
}
