package utils

import (
	"reflect"
	"testing"
)

func TestPermutations1(t *testing.T) {
	input := [][]string{
		{"1", "2"},
		{"3", "4"},
	}

	expected := [][]string{
		{"1", "3"},
		{"1", "4"},
		{"2", "3"},
		{"2", "4"},
	}

	actual := PermuteStrings(input...)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("permutations do not match")
	}
}
