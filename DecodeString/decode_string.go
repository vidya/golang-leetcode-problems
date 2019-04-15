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

//token stack
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
func tokenStream(str string) <-chan Token {
	tChan := make(chan Token)

	pos := 0
	go func() {
		strLen := len(str)
		var token Token

		for pos < strLen {
			switch rune(str[pos]) {
			case '[':
				token = Token{OPEN_BRACKET, "["}
				pos += 1

			case ']':
				token = Token{CLOSE_BRACKET, "["}
				pos += 1

			default:
				tKind := NUMBER
				if unicode.IsLetter(rune(str[pos])) {
					tKind = STRING
				}

				tStr := ""
				for n := 0; pos+n < strLen; n++ {
					ch := rune(str[pos+n])

					if (tKind == STRING) && !unicode.IsLetter(ch) {
						break
					}

					if (tKind == NUMBER) && !unicode.IsDigit(ch) {
						break
					}

					tStr += string(ch)
				}

				token = Token{tKind, tStr}
				pos += len(tStr)
			}

			tChan <- token
		}

		tChan <- Token{END_OF_STRING, "END_OF_STRING"}
		close(tChan)
	}()

	return tChan
}

// end: token stream

func DecodeString(encodedStr string) string {
	if encodedStr == "" {
		return ""
	}

	tStk := NewStack()
	tStk.Push(Token{STACK_BOTTOM, STACK_BOTTOM})

	for token := range tokenStream(encodedStr) {
		switch token.kind {
		case STRING:
			tv := token.value

			// if the stack top  is a STRING, combine this string with stack top
			tKind, tStr := tStk.Peek()
			if tKind == STRING {
				tv = tStr + tv

				tStk.Pop()
			}

			tStk.Push(Token{STRING, tv})

		case CLOSE_BRACKET:
			// Pop the STRING enclosed in brackets, OPEN_BRACKETand the preceding NUMBER
			strToken := tStk.Pop() // enclosed STRING
			tStk.Pop()             // OPEN_BRACKET
			numToken := tStk.Pop() // NUMBER

			repeatCount, _ := strconv.Atoi(numToken.value)

			tv := ""
			for n := 0; n < repeatCount; n++ {
				tv += strToken.value
			}

			// check if we have a STRING on top of stack
			tKind, tStr := tStk.Peek()
			if tKind == STRING {
				tv = tStr + tv

				tStk.Pop()
			}

			tStk.Push(Token{STRING, tv})

		case END_OF_STRING:
			break

		default:
			tStk.Push(token)
		}
	}

	_, tStr := tStk.Peek()
	return tStr
}
