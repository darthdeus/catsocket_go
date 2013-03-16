package main

import "testing"

func TestFoo(t *testing.T) {
	t.Errorf("John %v", true)
}
