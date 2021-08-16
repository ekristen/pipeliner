package utils

import "testing"

func TestStringSliceContainsPositive(t *testing.T) {
	test1 := []string{"one", "two", "three"}

	if !StringSliceContains(test1, "one") {
		t.Errorf("unable to find expected element in slice array")
	}
}

func TestStringSliceContainsNegative(t *testing.T) {
	test1 := []string{"one", "two", "three"}

	if StringSliceContains(test1, "five") {
		t.Errorf("found an element that should not be found")
	}
}

func TestStringSlicePosition1(t *testing.T) {
	test1 := []string{"one", "two", "three"}

	if pos := StringSlicePosition(test1, "one"); pos != 0 {
		t.Errorf("found incorrect position index (actual: %d)", pos)
	}
}

func TestStringSlicePosition2(t *testing.T) {
	test1 := []string{"one", "two", "three"}

	if pos := StringSlicePosition(test1, "three"); pos != 2 {
		t.Errorf("found incorrect position index (actual: %d)", pos)
	}
}

func TestStringSlicePosition3(t *testing.T) {
	test1 := []string{"one", "two", "three"}

	if pos := StringSlicePosition(test1, "four"); pos != -1 {
		t.Errorf("found incorrect position index (actual: %d)", pos)
	}
}
