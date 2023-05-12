package main

import (
	"fmt"
	. "jackcompiler/pkg/analyzer"
	. "jackcompiler/pkg/common"
)

func main() {
	tokenizer := NewTokenizer("tests/ArrayTest/Main.jack")

	for tokenizer.HasMoreTokens() {
		token := tokenizer.NextToken()
		switch token.TokenType() {
		case Keyword:
			fmt.Println(token.KeywordType())
		case Symbol:
			fmt.Println(string(token.Symbol()))
		case IntegerConstant:
			fmt.Println(token.IntVal())
		case StringConstant:
			fmt.Println(token.StringVal())
		case Identifier:
			fmt.Println(token.Identifier())
		}
	}

}
