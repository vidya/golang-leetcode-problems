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

package main

import (
	"strconv"
	"unicode"
)

const (
	END_OF_STRING = string(iota)
	STRING
	NUMBER
	OPEN_BRACKET
	CLOSE_BRACKET

	STACK_BOTTOM
)

// token stack
type tokenType struct {
	kind  string
	value string
}

var tokenStack []tokenType

func pushToken(newToken tokenType) {
	tokenStack = append(tokenStack, newToken)
}

func popToken() tokenType {
	tokenCount := len(tokenStack)
	topToken := tokenStack[tokenCount-1]

	tokenStack = tokenStack[0 : tokenCount-1]

	return topToken
}

func peekToken() (string, string) {
	tokenCount := len(tokenStack)
	topToken := tokenStack[tokenCount-1]

	return topToken.kind, topToken.value
}

// end: token stack

// token stream
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
	defer close(tokenChan)

	strLen := len(str)
	count := 0
	for count < strLen {
		switch ch := rune(str[count]); ch {
		case '[':
			tokenChan <- tokenType{OPEN_BRACKET, "["}
			count += 1

		case ']':
			tokenChan <- tokenType{CLOSE_BRACKET, "]"}
			count += 1

		default:
			token := getToken(str, count)
			tokenChan <- token

			count += len(token.value)
		}
	}

	tokenChan <- tokenType{END_OF_STRING, "END_OF_STRING"}
}

// end: token stream

func DecodeString(encodedStr string) string {
	if encodedStr == "" {
		return ""
	}

	tokenInputChan := make(chan tokenType, 2)

	go tokenStream(encodedStr, tokenInputChan)

	pushToken(tokenType{STACK_BOTTOM, STACK_BOTTOM})

	endOfString := false
	for !endOfString {
		switch token := <-tokenInputChan; token.kind {
		case STRING:
			lastStr := token.value

			// if the stack top  is a STRING, combine this string with stack top
			tokenKind, tokenValue := peekToken()
			if tokenKind == STRING {
				lastStr = tokenValue + lastStr

				popToken()
			}

			pushToken(tokenType{STRING, lastStr})

		case CLOSE_BRACKET:
			// pop the STRING enclosed in brackets, OPEN_BRACKETand the preceding NUMBER
			strToken := popToken() // enclosed STRING
			popToken()             // OPEN_BRACKET
			numToken := popToken() // NUMBER

			repeatCount, _ := strconv.Atoi(numToken.value)

			outStr := ""
			for n := 0; n < repeatCount; n++ {
				outStr += strToken.value
			}

			// check if we have a STRING on top of stack
			tokenKind, tokenValue := peekToken()
			if tokenKind == STRING {
				outStr = tokenValue + outStr

				popToken()
			}

			pushToken(tokenType{STRING, outStr})

		case END_OF_STRING:
			endOfString = true

		default:
			pushToken(token)
		}
	}

	_, tokenValue := peekToken()
	return tokenValue
}
