package utils

import "testing"

func TestRandomStringLength1(t *testing.T) {
	s := RandomString(12)

	if len(s) != 12 {
		t.Errorf("random string length")
	}
}

func TestRandomStringLength2(t *testing.T) {
	s := RandomString(128)

	if len(s) != 128 {
		t.Errorf("random string length")
	}
}

func TestRandomStringWithCharset1(t *testing.T) {
	s := RandomStringWithCharset(16, "a")

	if len(s) != 16 {
		t.Error("random string length is wrong")
	}

	if s != "aaaaaaaaaaaaaaaa" {
		t.Error("random string not generated properly")
	}
}
