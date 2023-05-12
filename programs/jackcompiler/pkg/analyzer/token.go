package analyzer

import (
	. "jackcompiler/pkg/common"
)

// Token is a struct that represents a token found within a .jack file
type Token struct {
	tokenType   TokenType
	keywordType KeywordType
	symbol      rune
	identifier  string
	intVal      int
	stringVal   string
}

// TokenType returns the type of the token
func (t *Token) TokenType() TokenType {
	return t.tokenType
}

// KeywordType returns the type of the keyword (if it is a keyword)
func (t *Token) KeywordType() KeywordType {
	return t.keywordType
}

// Symbol returns the symbol of the token (if it is a symbol)
func (t *Token) Symbol() rune {
	return t.symbol
}

// Identifier returns the identifier of the token (if it is an identifier)
func (t *Token) Identifier() string {
	return t.identifier
}

// IntVal returns the integer value of the token (if it is an integer constant)
func (t *Token) IntVal() int {
	return t.intVal
}

// StringVal returns the string value of the token (if it is a string constant)
func (t *Token) StringVal() string {
	return t.stringVal
}
