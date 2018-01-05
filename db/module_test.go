package db

import "testing"

func TestQuestionMarks(t *testing.T) {
	tests := []struct {
		n      int
		expect string
	}{
		{1, "(?)"},
		{2, "(?,?)"},
		{3, "(?,?,?)"},
	}

	for _, test := range tests {
		got := QuestionMarks(test.n)
		if got != test.expect {
			t.Fatalf("expected to get QuestionMarks(%d)=%s; got %s", test.n, test.expect, got)
		}
	}
}
