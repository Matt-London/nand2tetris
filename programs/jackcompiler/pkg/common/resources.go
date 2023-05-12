package common

import (
	"regexp"
)

// KeywordMap will map keywords by string to their respective keyword type
var KeywordMap = map[string]KeywordType{ // While not technically a constant
	"class":       Class,
	"method":      Method,
	"function":    Function,
	"constructor": Constructor,
	"int":         Int,
	"boolean":     Boolean,
	"char":        Char,
	"void":        Void,
	"var":         Var,
	"static":      Static,
	"field":       Field,
	"let":         Let,
	"do":          Do,
	"if":          If,
	"else":        Else,
	"while":       While,
	"return":      Return,
	"true":        True,
	"false":       False,
	"null":        Null,
	"this":        This,
}

// Regex expressions used by the tokenizer

var WhitespaceRegex = regexp.MustCompile(`^\s+`)
var EmptyRegex = regexp.MustCompile(`^\s*$`)
var CommentRegex = regexp.MustCompile(`^(//.*)|(/\*([^/]|$)|[^*])*\*/`)
var KeywordRegex = regexp.MustCompile(`^(class|constructor|function|method|field|static|var|int|char|boolean|void|true|false|null|this|let|do|if|else|while|return)\b`)
var SymbolRegex = regexp.MustCompile(`^[{}()\[\].,;+\-*/&|<>=~]`)
var IntegerConstantRegex = regexp.MustCompile(`^\d+`)
var StringConstantRegex = regexp.MustCompile(`^"([^"\n]*)"`)
var IdentifierRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`)
