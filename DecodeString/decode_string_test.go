package main

import (
	"fmt"
	"os"
	"testing"
)

type paramsType struct {
	testName   string
	encodedStr string
	want       string
}

func TestDecodeString(t *testing.T) {
	testFunc := func(params paramsType) {
		t.Helper()

		got := DecodeString(params.encodedStr)

		if got != params.want {
			fmt.Printf("\n FAIL: %s: (encodedStr, want,got) = (%s, %s, %s)\n",
				params.testName, params.encodedStr, params.want, got)

			t.Errorf("%s: got '%s' want '%s'", params.testName, got, params.want)
			os.Exit(3)
		}

		fmt.Printf("PASS: %s\n", params.testName)
	}

	testFunc(paramsType{"test_1_2[a]", "2[a]", "aa"})

	testFunc(paramsType{"test_2_3[a]2[bc]", "3[a]2[bc]", "aaabcbc"})

	testFunc(paramsType{"test_3_2[abc]3[cd]ef", "2[abc]3[cd]ef", "abcabccdcdcdef"})

	testFunc(paramsType{"test_4_12[x]", "12[x]", "xxxxxxxxxxxx"})

	testFunc(paramsType{"test_5_3[y]", "3[y]", "yyy"})

	testFunc(paramsType{"test_6_3[a2[c]]", "3[a2[c]]", "accaccacc"})

	testFunc(paramsType{"test_7_3[a2[c]]", "3[a]2[b4[F]c]", "aaabFFFFcbFFFFc"})
}
