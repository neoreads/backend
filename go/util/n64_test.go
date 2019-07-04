package util

import (
	"testing"
)

func TestN64(t *testing.T) {
	got := ParseN64("0")
	if got != 0 {
		t.Errorf("ParseN64(\"0\") = %d; want 0", got)
	}
}

func TestGenN64(t *testing.T) {
	got := GenN64(5)
	t.Logf("GenN64:%s", got)
	if len(got) != 5 {
		t.Errorf("len(TestGenN64(\"5\")) = %d; want 5", len(got))
	}
}
