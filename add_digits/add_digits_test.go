package main

import (
	"fmt"
	"os"
	"testing"
)

func testFunc(t *testing.T, test_name string, num int, want int) {
	got, want := AddDigits(num), want

	if got != want {
		t.Errorf("%s: got '%d' want '%d'", test_name, got, want)
		fmt.Printf("\nFAIL: %s\n\n", test_name)

		os.Exit(3)
	}

	fmt.Printf("PASS: %s\n", test_name)
}

func TestAddDigits(t *testing.T) {
	testFunc(t, "test_1_13", 13, 4)
	testFunc(t, "test_2_267", 267, 6)
	testFunc(t, "test_3_5013", 5013, 9)

	testFunc(t, "test_4_5", 5, 5)

}
