package util

import (
	"fmt"
	"strings"
)

func QuestionMarks(n int) string {
	if n <= -1 {
		panic(fmt.Errorf("programming error: %v as input for QuestionMarks()", n))
	}
	return fmt.Sprintf("(%s?)", strings.Repeat("?,", n-1))
}
