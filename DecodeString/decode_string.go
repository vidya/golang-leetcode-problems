Last login: Sun Feb 10 13:47:41 on ttys023
>pwd
/Users/vidya
>cd Documents/Project/*Go/Git*
>cd golang-leetcode-problems
>git status
On branch master
Your branch is up to date with 'origin/master'.

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	DecodeString/

nothing added to commit but untracked files present (use "git add" to track)
>git add DecodeString/
>git status
On branch master
Your branch is up to date with 'origin/master'.

Changes to be committed:
  (use "git reset HEAD <file>..." to unstage)

	new file:   DecodeString/decode_string.go
	new file:   DecodeString/decode_string_test.go

>git commit -m "solution to DecodeString problem" -a
[master 4f295cd] solution to DecodeString problem
 2 files changed, 225 insertions(+)
 create mode 100644 DecodeString/decode_string.go
/*
        https://leetcode.com/problems/decode-string/

        394. Decode String

        Given an encoded string, return it's decoded string.

        The encoding rule is: k[encoded_string], where the encoded_string
        inside the square brackets is being repeated exactly k times.
        Note that k is guaranteed to be a positive integer.

        You may assume that the input string is always valid; No extra
        white spaces, square brackets are well-formed, etc.

        Furthermore, you may assume that the original data does not
        contain any digits and that digits are only for those repeat
        numbers, k. For example, there won't be input like 3a or 2[4].

        Examples:
        s = "3[a]2[bc]", return "aaabcbc".
        s = "3[a2[c]]", return "accaccacc".
        s = "2[abc]3[cd]ef", return "abcabccdcdcdef".
*/

"DecodeString/decode_string.go" 180L, 3776C/*
	https://leetcode.com/problems/decode-string/

	394. Decode String

	Given an encoded string, return it's decoded string.

	The encoding rule is: k[encoded_string], where the encoded_string
	inside the square brackets is being repeated exactly k times.
	Note that k is guaranteed to be a positive integer.

	You may assume that the input string is always valid; No extra
	white spaces, square brackets are well-formed, etc.

	Furthermore, you may assume that the original data does not
	contain any digits and that digits are only for those repeat
	numbers, k. For example, there won't be input like 3a or 2[4].

	Examples:
	s = "3[a]2[bc]", return "aaabcbc".
	s = "3[a2[c]]", return "accaccacc".
	s = "2[abc]3[cd]ef", return "abcabccdcdcdef".
*/

package main

import (
	"fmt"
	"strconv"
	"unicode"
)

const ENDOFSTRING = "ENDOFSTRING"
const STRING = "STRING"
const NUMBER = "NUMBER"
const OPENBRACKET = "OPENBRACKET"
const CLOSEBRACKET = "CLOSEBRACKET"

const STACKBOTTOM = "STACKBOTTOM"

type tokenType struct {
	kind  string
	value string
}

var tokenStack []tokenType

func getToken(str string, pos int) tokenType {
	tokenKind := NUMBER
	if unicode.IsLetter(rune(str[pos])) {
		tokenKind = STRING
	}

	tokenValue := ""
	strLen := len(str)
	for n := 0; pos+n < strLen; n++ {
		ch := rune(str[pos+n])

		if (tokenKind == STRING) && !unicode.IsLetter(ch) {
			break
		}

		if (tokenKind == NUMBER) && !unicode.IsDigit(ch) {
			break
		}

		tokenValue += string(ch)
	}

	return tokenType{tokenKind, tokenValue}
}

func tokenStream(str string, tokenChan chan tokenType) {
	fmt.Printf("tokenStream(): encodedStr  = %s\n", str)
	strLen := len(str)

	count := 0
	for count < strLen {
		switch ch := rune(str[count]); ch {
		case '[':
			tokenChan <- tokenType{OPENBRACKET, "["}
			count += 1

		case ']':
			tokenChan <- tokenType{CLOSEBRACKET, "]"}
			count += 1

		default:
			token := getToken(str, count)
			tokenChan <- token

			count += len(token.value)
		}
	}

	tokenChan <- tokenType{ENDOFSTRING, "ENDOFSTRING"}
	close(tokenChan)
}

func pushToken(newToken tokenType) {
	tokenStack = append(tokenStack, newToken)
}

func popTokenStack() tokenType {
	tokenCount := len(tokenStack)
	topToken := tokenStack[tokenCount-1]

	tokenStack = tokenStack[0 : tokenCount-1]

	return topToken
}

func getTopToken() (string, string) {
	tokenCount := len(tokenStack)
	topToken := tokenStack[tokenCount-1]

	return topToken.kind, topToken.value
}

func DecodeString(encodedStr string) string {
	if encodedStr == "" {
		return ""
	}

	tokenInputChan := make(chan tokenType, 2)

	go tokenStream(encodedStr, tokenInputChan)

	pushToken(tokenType{STACKBOTTOM, STACKBOTTOM})

	endOfString := false
	for !endOfString {
		switch token := <-tokenInputChan; token.kind {
		case STRING:
			lastStr := token.value

			// if the stack top  is a STRING, combine this string with stack top
			tokenKind, tokenValue := getTopToken()
			if tokenKind == STRING {
				lastStr = tokenValue + lastStr

				popTokenStack()
			}

			pushToken(tokenType{STRING, lastStr})

		case CLOSEBRACKET:
			// pop the STRING enclosed in brackets, OPENBRACKET and the preceding NUMBER
			strToken := popTokenStack() // enclosed STRING
			popTokenStack()             // OPENBRACKET
			numToken := popTokenStack() // NUMBER

			repeatCount, _ := strconv.Atoi(numToken.value)

			outStr := ""
			for n := 0; n < repeatCount; n++ {
				outStr += strToken.value
			}

			// check if we have a STRING on top of stack
			tokenKind, tokenValue := getTopToken()
			if tokenKind == STRING {
				outStr = tokenValue + outStr

				popTokenStack()
			}

			pushToken(tokenType{STRING, outStr})

		case ENDOFSTRING:
			endOfString = true

		default:
			pushToken(token)
		}
	}

	_, tokenValue := getTopToken()
	return tokenValue
}
