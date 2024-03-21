package main

import (
	"unicode"
)

type TokenType uint
type Token struct {
	Type TokenType
	Line int
	Row  int
	Data string
}

const (
	NameToken TokenType = iota
	OpenBracketToken
	CloseBracketToken
	CommaToken
)

func TokenTypeName(tokenType TokenType) string {
	switch tokenType {
	case NameToken:
		return "NameToken"
	case OpenBracketToken:
		return "OpenBracketToken"
	case CloseBracketToken:
		return "CloseBracketToken"
	case CommaToken:
		return "CommaToken"
	}

	return "unknown"
}

func Tokenize(text string) []*Token {
	var tokens []*Token
	var lastToken *Token

	line := 0
	lineStart := 0

	for i, char := range text {
		switch {

		case char == '(':
			tokens = append(tokens, &Token{
				Type: OpenBracketToken,
				Line: line,
				Row:  i - lineStart,
				Data: "(",
			})
			lastToken = nil
		case char == ')':
			tokens = append(tokens, &Token{
				Type: CloseBracketToken,
				Line: line,
				Row:  i - lineStart,
				Data: ")",
			})
			lastToken = nil

		case char == ',':
			tokens = append(tokens, &Token{
				Type: CommaToken,
				Line: line,
				Row:  i - lineStart,
				Data: ",",
			})
			lastToken = nil

		case char == ' ':
			lastToken = nil
		case char == '\n':
			line++
			lineStart = i + 1
			lastToken = nil

		case unicode.IsLetter(char):
			if lastToken == nil {
				lastToken = &Token{
					Type: NameToken,
					Line: line,
					Row:  i - lineStart,
					Data: "",
				}

				tokens = append(tokens, lastToken)
			}

			lastToken.Data += string(char)
		}

	}

	return tokens
}
