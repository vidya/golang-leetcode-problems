package main

import (
	"fmt"
	"os"
	"testing"
)

type paramsType struct {
	testName string
	num      int
	want     int
}


func TestAddDigits(t *testing.T) {
	testFunc := func(params paramsType) {
		t.Helper()

		got := AddDigits(params.num)

		if got != params.want {
			t.Errorf("%s: got '%d' want '%d'", params.testName, got, params.want)
			fmt.Printf("\nFAIL: %s\n\n", params.testName)

			os.Exit(3)
		}

		fmt.Printf("PASS: %s\n", params.testName)
	}

	testFunc(paramsType{"test_1_13", 13, 4})
	testFunc(paramsType{"test_2_267", 267, 6})
	testFunc(paramsType{"test_3_5013", 5013, 9})
	testFunc(paramsType{"test_4_5", 5, 5})
}
