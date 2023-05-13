package analyzer

import (
	"bufio"
	"fmt"
	. "jackcompiler/pkg/common"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Tokenizer takes in a file and returns tokens as needed
type Tokenizer struct {
	inputText    string
	prevMatchEnd int
}

// NewTokenizer takes the input file path and loads a new tokenizer
func NewTokenizer(inputFilePath string) *Tokenizer {
	// Read the file
	file, err := os.Open(inputFilePath)
	if err != nil {
		panic(err)
		return nil
	}

	// Close the file when we are done, panic if we hit an error
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var contents string
	for scanner.Scan() {
		contents += scanner.Text() + "\n"
	}

	return &Tokenizer{inputText: contents, prevMatchEnd: -1}

}

// matchToken will return the current token
// This assumes that the regex matches the beginning of the inputText
func (t *Tokenizer) matchToken(regex regexp.Regexp) string {
	// Get the indexes of the match
	matchEndIdx := regex.FindStringIndex(t.inputText)[1]

	// Collect string
	match := t.inputText[:matchEndIdx]
	// Get ready to advance the inputText when asked
	t.prevMatchEnd = matchEndIdx

	return match
}

// HasMoreTokens returns whether there are more tokens in the input
func (t *Tokenizer) HasMoreTokens() bool {
	// Check if whitespace regex matches the entire remaining string
	return !EmptyRegex.MatchString(t.inputText)
}

// Advance will advance the tokenizer to the next token
func (t *Tokenizer) Advance() {
	// Advance the input text
	if t.prevMatchEnd != -1 {
		t.inputText = t.inputText[t.prevMatchEnd:]
		t.prevMatchEnd = -1
	}
}

// matchAdvance will match the regex and advance the tokenizer
func (t *Tokenizer) matchAdvance(regex regexp.Regexp) string {
	match := t.matchToken(regex)
	t.Advance()
	return match
}

// Token returns the current token
func (t *Tokenizer) Token() *Token {
	token := &Token{}

	// Throwaway whitespace matching whitespace regex
	// And throwaway comment regex
	changed := true
	for changed {
		changed = false
		// Clear whitespace
		if WhitespaceRegex.MatchString(t.inputText) {
			t.matchAdvance(*WhitespaceRegex)
			changed = true
		}
		// Clear comments
		if CommentRegex.MatchString(t.inputText) {
			t.matchAdvance(*CommentRegex)
			changed = true
		}

	}

	// We need to check in order of precedence, for example checking an identifier before a symbol would resolve
	// "int" as a symbol, which is not right
	if t.HasMoreTokens() {
		if KeywordRegex.MatchString(t.inputText) {
			// Next token is keyword
			token.tokenType = Keyword

			match := t.matchToken(*KeywordRegex)

			// Check which keyword it is
			token.keywordType = KeywordMap[match]

		} else if SymbolRegex.MatchString(t.inputText) {
			// Next token is a symbol
			token.tokenType = Symbol

			match := t.matchToken(*SymbolRegex)

			// Symbol should just be a char
			token.symbol = rune(match[0])

		} else if IntegerConstantRegex.MatchString(t.inputText) {
			// Next token is an integer constant
			token.tokenType = IntegerConstant

			match := t.matchToken(*IntegerConstantRegex)

			// Convert string to int
			token.intVal, _ = strconv.Atoi(match)

		} else if StringConstantRegex.MatchString(t.inputText) {
			// Next token is a string constant
			token.tokenType = StringConstant

			match := t.matchToken(*StringConstantRegex)

			// Remove quotes from string and set its value
			token.stringVal = match[1 : len(match)-1]

		} else if IdentifierRegex.MatchString(t.inputText) {
			// Next token is an identifier
			token.tokenType = Identifier

			match := t.matchToken(*IdentifierRegex)

			// Set identifier value
			token.identifier = match
		} else {
			// Next token is invalid
			fmt.Println("Unexpected token: " + strings.Split(t.inputText, " ")[0])
		}
		return token
	}
	// If it doesn't have more tokens, then return nil, but this should never happen
	return nil
}
