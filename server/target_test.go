package server

import "testing"

func TestTargetHandler_extractName(t *testing.T) {
	h := &TargetHandler{}
	tests := []struct {
		title     string
		ok        bool
		name      string
		condition string
	}{
		{
			title: "Motoracer 4 - PlayStation VR",
			ok:    false,
		},
		{
			title: "RiME - PlayStation 4",
			ok:    true,
			name:  "RiME",
		},
		{
			title: "NBA Live 16 (PlayStation 4)",
			ok:    true,
			name:  "NBA Live 16",
		},
		{
			title: "Uncharted 4: A Thief&rsquo;s End&#153; Special Edition (PlayStation 4)",
			ok:    true,
			name:  "Uncharted 4: A Thief’s End™ Special Edition",
		},
		{
			title: "Dragonball XenoversePlaystation 4",
			ok: true,
			name: "Dragonball Xenoverse",
		},
	}
	for _, test := range tests {
		nameGot, okGot := h.extractName(test.title)
		if test.ok && !okGot {
			t.Fatalf("expected to parse %v successfully", test.title)
		}
		if !test.ok {
			continue
		}
		if test.name != nameGot {
			t.Fatalf("expected to get name %v; got %v", test.name, nameGot)
		}
	}
}
