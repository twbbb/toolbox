// Package arrayutil arrayutil
package arrayutil

import (
	"testing"
)

// Test_ReverseSlice
func Test_ReverseSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	ReverseSlice(slice)
	if slice[0] != 5 {
		t.Errorf("ReverseSlice failed")
	}

}
