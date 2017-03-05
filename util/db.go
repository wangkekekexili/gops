package util

import (
	"fmt"
	"strings"
)

// QuestionMarks generates something like (?,?,?) to be used by SQL statement.
// Caller must guarantee that input n is positive.
func QuestionMarks(n int) string {
	if n <= 0 {
		panic(fmt.Errorf("programming error: %v as input for QuestionMarks()", n))
	}
	return fmt.Sprintf("(%s?)", strings.Repeat("?,", n-1))
}
