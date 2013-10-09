package broadway

import (
  "testing"
  "strings"
)

func TestConcat(t *testing.T) {
  a := []string{"A"}
  b := []string{"B"}
  c := concat(a, b)
  result := strings.Join(c, ", ")
  expected := "A, B"
  if result != expected {
    t.Errorf("Expected %v to be %v", result, expected)
  }
}

