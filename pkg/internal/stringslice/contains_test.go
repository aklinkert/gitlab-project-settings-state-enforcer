package stringslice

import "testing"

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !Contains("b", slice) {
		t.Errorf("Expected contains to returns true as the slice contains the element, but it returned false")
	}

	if Contains("x", slice) {
		t.Errorf("Expected contains to returns false as the slice does not contain the element, but it returned true")
	}
}
