package common

// TokenType is an enum for the type of token
type TokenType int

const (
	Keyword TokenType = iota
	Symbol
	IntegerConstant
	StringConstant
	Identifier
)

// KeywordType is an enum for type of keyword a token has (if it is a keyword)
type KeywordType int

const (
	Class KeywordType = iota
	Method
	Function
	Constructor
	Int
	Boolean
	Char
	Void
	Var
	Static
	Field
	Let
	Do
	If
	Else
	While
	Return
	True
	False
	Null
	This
)
