package oth

import (
	"fmt"
	"testing"
)

func TestStream(t *testing.T) {
	var s Optional[string]
	s.OfNilable([]string{"go", "java", "python"}).IfPresent(func(s string) {
		fmt.Println(s)
	})
}
