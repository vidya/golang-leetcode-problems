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
type Token struct {
	kind  string
	value string
}

type Stack struct {
	stack []Token
}

type stackInterface interface {
	Push(token Token)
	Pop() Token
	Peek() (string, string)
	Len() int
	NewStack() *Stack
}

func (stk *Stack) Push(token Token) {
	stk.stack = append(stk.stack, token)
}

func (stk *Stack) Pop() Token {
	count := stk.Len()

	top := stk.stack[count-1]
	stk.stack = stk.stack[:count-1]

	return top
}

func (stk *Stack) Peek() (string, string) {
	count := stk.Len()
	top := stk.stack[count-1]

	return top.kind, top.value
}

func (stk *Stack) Len() int {
	return len(stk.stack)
}

func NewStack() *Stack {
	return &Stack{}
}
// end: token stack

// token stream
func getToken(str string, pos int) Token {
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

	return Token{tokenKind, tokenValue}
}

func tokenStream(str string, tokenChan chan Token) {
	defer close(tokenChan)

	strLen := len(str)
	count := 0
	for count < strLen {
		switch ch := rune(str[count]); ch {
		case '[':
			tokenChan <- Token{OPEN_BRACKET, "["}
			count++

		case ']':
			tokenChan <- Token{CLOSE_BRACKET, "]"}
			count++

		default:
			token := getToken(str, count)
			tokenChan <- token

			count += len(token.value)
		}
	}

	tokenChan <- Token{END_OF_STRING, "END_OF_STRING"}
}
// end: token stream

func DecodeString(encodedStr string) string {
	if encodedStr == "" {
		return ""
	}

	tokenStack := NewStack()

	tokenInputChan := make(chan Token, 2)

	go tokenStream(encodedStr, tokenInputChan)

	tokenStack.Push(Token{STACK_BOTTOM, STACK_BOTTOM})

	endOfString := false
	for !endOfString {
		switch token := <-tokenInputChan; token.kind {
		case STRING:
			lastStr := token.value

			// if the stack top  is a STRING, combine this string with stack top
			tokenKind, tokenValue := tokenStack.Peek()
			if tokenKind == STRING {
				lastStr = tokenValue + lastStr

				tokenStack.Pop()
			}

			tokenStack.Push(Token{STRING, lastStr})

		case CLOSE_BRACKET:
			// Pop the STRING enclosed in brackets, OPEN_BRACKETand the preceding NUMBER
			strToken := tokenStack.Pop() // enclosed STRING
			tokenStack.Pop()             // OPEN_BRACKET
			numToken := tokenStack.Pop() // NUMBER

			repeatCount, _ := strconv.Atoi(numToken.value)

			outStr := ""
			for n := 0; n < repeatCount; n++ {
				outStr += strToken.value
			}

			// check if we have a STRING on top of stack
			tokenKind, tokenValue := tokenStack.Peek()
			if tokenKind == STRING {
				outStr = tokenValue + outStr

				tokenStack.Pop()
			}

			tokenStack.Push(Token{STRING, outStr})

		case END_OF_STRING:
			endOfString = true

		default:
			tokenStack.Push(token)
		}
	}

	_, tokenValue := tokenStack.Peek()
	return tokenValue
}
