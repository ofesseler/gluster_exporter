package main

import "testing"

func TestContainsVolume(t *testing.T) {
	expamle := "doge"
	testSlice := []string{"wow", "such", expamle}
	if !ContainsVolume(testSlice, expamle) {
		t.Fatalf("Hasn't found %v in slice %v", expamle, testSlice)
	}
}
