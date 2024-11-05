package test

import (
	"testing"
)

const ret string = "Helo"

func TestAssertHello(t *testing.T) {
	if ret != "Hello" {
		t.Fatalf("should have been test")
	}
}
