package analyzer

import (
	"fmt"
	. "jackcompiler/pkg/common"
	"os"
	"strconv"
	"strings"
)

// Engine is the handler of compilation
type Engine struct {
	tokenizer  *Tokenizer
	indentVal  int
	inputPath  string
	outputFile *os.File
}

// NewEngine constructs an engine and tokenizer from an input file
func NewEngine(inputFilePath string) *Engine {
	return &Engine{tokenizer: NewTokenizer(inputFilePath), inputPath: inputFilePath}
}

// write Will write a string to the output file, obeying indentation and adding new lines
// The caller is responsible for setting the indentation value
func (e *Engine) write(strToWrite string) {
	formattedOutput := strings.Repeat("\t", e.indentVal) + strToWrite + "\n"
	_, err := e.outputFile.WriteString(formattedOutput)
	if err != nil {
		panic(err)
		return
	}
}

// writeKeyword will write a keyword's xml to the output file, then advance
// It will return false if the keyword is not the next token
func (e *Engine) writeKeyword() bool {
	if !(e.tokenizer.Token().TokenType() == Keyword) {
		return false
	}
	e.write("<keyword> " + KeywordStrMap[e.tokenizer.Token().KeywordType()] + " </keyword>")
	e.tokenizer.Advance()

	return true
}

// writeSpecSymbol will write a specific symbol's xml to the output file, then advance
// It will return false if the symbol is not the next token
func (e *Engine) writeSpecSymbol(symbol rune) bool {
	if !(e.tokenizer.Token().TokenType() == Symbol && e.tokenizer.Token().Symbol() == symbol) {
		return false
	}
	e.writeSymbol()
	return true
}

// writeType will write a type's xml to the output file, then advance
// It will return false if the type is not the next token
func (e *Engine) writeType() bool {
	// Check if we have an identifier
	if e.tokenizer.Token().TokenType() == Identifier {
		// This can be a type
		e.writeIdentifier()
		return true
	} else if e.tokenizer.Token().TokenType() == Keyword {
		// Make sure this is either int, char, or boolean
		if e.tokenizer.Token().KeywordType() == Int || e.tokenizer.Token().KeywordType() == Char ||
			e.tokenizer.Token().KeywordType() == Boolean || e.tokenizer.Token().KeywordType() == Void {
			e.writeKeyword()
			return true
		}
	}

	// Otherwise, we don't have a type
	return false
}

// writeSymbol will write a symbol's xml to the output file, then advance
// It will return false if the symbol is not the next token
func (e *Engine) writeSymbol() bool {
	if !(e.tokenizer.Token().TokenType() == Symbol) {
		return false
	}

	// Escape special characters
	cleanedSymbol := string(e.tokenizer.Token().Symbol())
	if cleanedSymbol == "<" {
		cleanedSymbol = "&lt;"
	} else if cleanedSymbol == ">" {
		cleanedSymbol = "&gt;"
	} else if cleanedSymbol == "&" {
		cleanedSymbol = "&amp;"
	}
	e.write("<symbol> " + cleanedSymbol + " </symbol>")
	e.tokenizer.Advance()

	return true
}

// writeIdentifier will write an identifier's xml to the output file, then advance
// It will return false if the identifier is not the next token
func (e *Engine) writeIdentifier() bool {
	if !(e.tokenizer.Token().TokenType() == Identifier) {
		return false
	}
	e.write("<identifier> " + e.tokenizer.Token().Identifier() + " </identifier>")
	e.tokenizer.Advance()

	return true
}

// compileClass will compile and write the xml for a class
// This will also handle the top level of the program
func (e *Engine) compileClass() bool {
	// Exit if we are trying to compile a class when no class keyword is available
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Class) {
		fmt.Println("Expected class keyword")
		return false
	}
	e.write("<class>")
	e.indentVal++

	// Now set keyword to class
	e.writeKeyword()

	// Now we should have a class name
	if !e.writeIdentifier() {
		fmt.Println("Expected identifier")
		return false
	}

	// Write the open brace
	e.writeSpecSymbol('{')

	// Now we move on to class variable declarations
	for e.compileClassVarDec() {
		// This will loop until we are done with class variable declarations
	}

	// Now we move on to subroutine declarations
	for e.compileSubroutine() {
		// This will loop until we are done with subroutine declarations
	}

	// Ensure that we have a closing brace
	e.writeSpecSymbol('}')

	e.indentVal--

	e.write("</class>")

	return true

}

// compileClassVarDec will write the xml for static and field variable declarations
func (e *Engine) compileClassVarDec() bool {
	// Next should be static or field
	if !(e.tokenizer.Token().TokenType() == Keyword && (e.tokenizer.Token().KeywordType() == Static || e.tokenizer.Token().KeywordType() == Field)) {
		return false
	}

	// Otherwise we do
	e.write("<classVarDec>")
	e.indentVal++

	// Grab the keyword (static or field)
	e.writeKeyword()

	// Now we should have a name of a type (int, char, boolean, or identifier)
	if !e.writeType() {
		fmt.Println("Expected type")
		return false
	}

	// Now we should have an identifier, we may have multiple so run this in a loop
	moreIdent := true
	for moreIdent {
		// Try to write an identifier
		if !e.writeIdentifier() {
			fmt.Println("Expected identifier")
			return false
		}

		// Check if we have a comma
		if e.tokenizer.Token().TokenType() == Symbol && e.tokenizer.Token().Symbol() == ',' {
			// Write the comma
			e.writeSymbol()
		} else {
			// Otherwise, we are done
			moreIdent = false
		}
	}

	// Now we should have a semicolon
	e.writeSpecSymbol(';')

	// Write the closing
	e.indentVal--
	e.write("</classVarDec>")

	return true

}

// compileSubroutine will write the xml for a subroutine
func (e *Engine) compileSubroutine() bool {
	// Check if next we have a constructor, function, or method
	if !(e.tokenizer.Token().TokenType() == Keyword && (e.tokenizer.Token().KeywordType() == Constructor ||
		e.tokenizer.Token().KeywordType() == Function || e.tokenizer.Token().KeywordType() == Method)) {
		// If not then we return false
		return false
	}

	// Write subroutine
	e.write("<subroutineDec>")
	e.indentVal++

	// Eat the function/constructor/method keyword
	e.writeKeyword()

	// Now we should have a return type
	if !e.writeType() {
		fmt.Println("Expected type")
		return false
	}

	// Next should be an identifier
	if !e.writeIdentifier() {
		fmt.Println("Expected identifier")
		return false
	}

	// Next should be an open parenthesis
	e.writeSpecSymbol('(')

	// Now we should have a parameter list
	// We don't need an if statement since it always prints the tags
	e.compileParameterList()

	// Now we should have a closing parenthesis
	e.writeSpecSymbol(')')

	// Now we should have a subroutine body
	e.write("<subroutineBody>")
	e.indentVal++

	// Now we should have an open brace
	e.writeSpecSymbol('{')

	// Now we should have variable declarations
	for e.compileVarDec() {
		// This will loop until we are done with variable declarations
	}

	// Now we should have statements
	for !((e.tokenizer.Token().TokenType() == Symbol) && (e.tokenizer.Token().Symbol() == '}')) {
		// This will loop until we are done with statements
		e.compileStatements()
	}

	// Now we should have a closing brace
	e.writeSpecSymbol('}')

	e.indentVal--
	e.write("</subroutineBody>")

	e.indentVal--
	e.write("</subroutineDec>")

	return true
}

// compileParameterList will write the xml for a parameter list
func (e *Engine) compileParameterList() bool {
	// Write opening
	e.write("<parameterList>")
	e.indentVal++

	moreParams := true

	// If the next token is a closing parenthesis then we are done
	if e.tokenizer.Token().TokenType() == Symbol && e.tokenizer.Token().Symbol() == ')' {
		moreParams = false
	}

	for moreParams {
		// Check if we have a type
		if !e.writeType() {
			fmt.Println("Expected type")
			return false
		}

		// Now we should have an identifier
		if !e.writeIdentifier() {
			fmt.Println("Expected identifier")
			return false
		}

		// Check if we have a comma
		if !e.writeSpecSymbol(',') {
			moreParams = false
		}
	}

	e.indentVal--
	e.write("</parameterList>")

	return true
}

// compileVarDec will write the xml for a variable declaration
func (e *Engine) compileVarDec() bool {
	// Check if we have a var keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Var) {
		return false
	}

	// Write opening
	e.write("<varDec>")
	e.indentVal++

	// Write var keyword
	e.writeKeyword()

	// Now we should have a type
	if !e.writeType() {
		fmt.Println("Expected type")
		return false
	}

	// We can have commas that cause multiple declarations so this is looped
	moreIdent := true
	for moreIdent {
		// Now we should have an identifier
		if !e.writeIdentifier() {
			fmt.Println("Expected identifier")
			return false
		}

		// Check if we have a comma
		if !e.writeSpecSymbol(',') {
			// Then we should write a semicolon
			if !e.writeSpecSymbol(';') {
				fmt.Println("Expected semicolon")
				return false
			}
			moreIdent = false
		}
	}

	e.indentVal--
	e.write("</varDec>")

	return true

}

// compileStatements will write the xml for a statement
func (e *Engine) compileStatements() bool {
	// This will write statements regardless
	e.write("<statements>")
	e.indentVal++

	// Loop while we have statements
	moreStatements := true
	for moreStatements {
		// Check if we have a let statement
		if e.compileLet() {
			continue
		}

		// Check if we have an if statement
		if e.compileIf() {
			continue
		}

		// Check if we have a while statement
		if e.compileWhile() {
			continue
		}

		// Check if we have a do statement
		if e.compileDo() {
			continue
		}

		// Check if we have a return statement
		if e.compileReturn() {
			continue
		}

		// If we get here then we don't have a statement
		moreStatements = false
	}

	e.indentVal--
	e.write("</statements>")

	return true
}

// compileLet will write the xml for a let statement
func (e *Engine) compileLet() bool {
	// Check if we have a let keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Let) {
		return false
	}

	// Write opening
	e.write("<letStatement>")
	e.indentVal++

	// Eat the keyword
	e.writeKeyword()

	// Now we should have an identifier
	if !e.writeIdentifier() {
		fmt.Println("Expected identifier")
		return false
	}

	// Next we may have an open brace or an equal sign
	if e.writeSpecSymbol('[') {
		// We have an open brace so we should have an expression
		if !e.compileExpression() {
			fmt.Println("Expected expression")
			return false
		}

		// Now we should have a closing brace
		if !e.writeSpecSymbol(']') {
			fmt.Println("Expected closing brace")
			return false
		}
	}

	// Now we should have an equal sign
	if !e.writeSpecSymbol('=') {
		fmt.Println("Expected equal sign")
		return false
	}

	// Now we should have an expression
	if !e.compileExpression() {
		fmt.Println("Expected expression")
		return false
	}

	// Now we should have a semicolon
	if !e.writeSpecSymbol(';') {
		fmt.Println("Expected semicolon")
		return false
	}

	e.indentVal--
	e.write("</letStatement>")

	return true
}

// compileIf will write the xml for an if statement
func (e *Engine) compileIf() bool {
	// Check if we have an if keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == If) {
		return false
	}

	// Write opening
	e.write("<ifStatement>")
	e.indentVal++

	// Eat if keyword
	e.writeKeyword()

	// Now we should have an open parenthesis
	if !e.writeSpecSymbol('(') {
		fmt.Println("Expected open parenthesis")
		return false
	}

	// Now we should have an expression
	if !e.compileExpression() {
		fmt.Println("Expected expression")
		return false
	}

	// Now we should have a close parenthesis
	if !e.writeSpecSymbol(')') {
		fmt.Println("Expected close parenthesis")
		return false
	}

	// Now we should have an open brace
	if !e.writeSpecSymbol('{') {
		fmt.Println("Expected open brace")
		return false
	}

	// Now we should have statements
	if !e.compileStatements() {
		fmt.Println("Expected statements")
		return false
	}

	// Now we should have a close brace
	if !e.writeSpecSymbol('}') {
		fmt.Println("Expected close brace")
		return false
	}

	// Now we may have an else keyword
	if e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Else {
		// Write the else keyword
		e.writeKeyword()

		// Now we should have an open brace
		if !e.writeSpecSymbol('{') {
			fmt.Println("Expected open brace")
			return false
		}

		// Now we should have statements
		if !e.compileStatements() {
			fmt.Println("Expected statements")
			return false
		}

		// Now we should have a close brace
		if !e.writeSpecSymbol('}') {
			fmt.Println("Expected close brace")
			return false
		}
	}

	e.indentVal--
	e.write("</ifStatement>")

	return true
}

// compileWhile will write the xml for a while statement
func (e *Engine) compileWhile() bool {
	// Check if we have a while keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == While) {
		return false
	}

	// Write opening
	e.write("<whileStatement>")
	e.indentVal++

	// Eat the keyword
	e.writeKeyword()

	// Now we should have an open parenthesis
	if !e.writeSpecSymbol('(') {
		fmt.Println("Expected open parenthesis")
		return false
	}

	// Now we should have an expression
	if !e.compileExpression() {
		fmt.Println("Expected expression")
		return false
	}

	// Now we should have a close parenthesis
	if !e.writeSpecSymbol(')') {
		fmt.Println("Expected close parenthesis")
		return false
	}

	// Now we should have an open brace
	if !e.writeSpecSymbol('{') {
		fmt.Println("Expected open brace")
		return false
	}

	// Now we should have statements
	if !e.compileStatements() {
		fmt.Println("Expected statements")
		return false
	}

	// Now we should have a close brace
	if !e.writeSpecSymbol('}') {
		fmt.Println("Expected close brace")
		return false
	}

	e.indentVal--
	e.write("</whileStatement>")

	return true

}

// compileDo will write the xml for a do statement
func (e *Engine) compileDo() bool {
	// We should have a do keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Do) {
		return false
	}

	// Write opening
	e.write("<doStatement>")
	e.indentVal++

	// Eat the keyword
	e.writeKeyword()

	// Now we should have an identifier
	if !e.writeIdentifier() {
		fmt.Println("Expected identifier")
		return false
	}

	// Now we may have a period
	if e.writeSpecSymbol('.') {
		// We have a period so we should have an identifier
		if !e.writeIdentifier() {
			fmt.Println("Expected identifier")
			return false
		}
	}

	// Now we should have an open parenthesis
	if e.writeSpecSymbol('(') {
		// We have an open parenthesis so we should have an expression list
		if !e.compileExpressionList() {
			fmt.Println("Expected expression list")
			return false
		}

		// Now we should have a closing parenthesis
		if !e.writeSpecSymbol(')') {
			fmt.Println("Expected closing parenthesis")
			return false
		}
	}

	// Now we should have a semicolon
	if !e.writeSpecSymbol(';') {
		fmt.Println("Expected semicolon")
		return false
	}

	e.indentVal--
	e.write("</doStatement>")

	return true
}

// compileReturn will write the xml for a return statement
func (e *Engine) compileReturn() bool {
	// We should have a return keyword
	if !(e.tokenizer.Token().TokenType() == Keyword && e.tokenizer.Token().KeywordType() == Return) {
		return false
	}

	// Write opening
	e.write("<returnStatement>")
	e.indentVal++

	// Eat return keyword
	e.writeKeyword()

	// We will have an expression or a semi colon
	statementEnded := e.tokenizer.Token().TokenType() == Symbol && e.tokenizer.Token().Symbol() == ';'

	// If we have an expression
	if !statementEnded {
		// We should have an expression
		if !e.compileExpression() {
			fmt.Println("Expected expression")
			return false
		}
	}

	// Now we should have a semicolon
	if !e.writeSpecSymbol(';') {
		fmt.Println("Expected semicolon")
		return false
	}

	e.indentVal--
	e.write("</returnStatement>")

	return true
}

// compileExpression will write the xml for an expression
func (e *Engine) compileExpression() bool {
	// Write opening
	e.write("<expression>")
	e.indentVal++

	// Compile terms until we don't have an operator
	moreTerms := true
	for moreTerms {
		if !e.compileTerm() {
			fmt.Println("Expected term")
			return false
		}

		// Check if we have an operator
		if e.tokenizer.Token().IsOperator() {
			// Write the operator
			e.writeSymbol()
		} else {
			moreTerms = false
		}
	}

	e.indentVal--
	e.write("</expression>")

	return true

}

// compileTerm will write the xml for a term (sorting out arrays vs calls and such)
func (e *Engine) compileTerm() bool {
	// A term can be an integer constant, string constant, keyword constant, variable name,
	// array, subroutine call, unary operation
	isInt := e.tokenizer.Token().TokenType() == IntegerConstant
	isString := e.tokenizer.Token().TokenType() == StringConstant
	isKeyword := e.tokenizer.Token().TokenType() == Keyword
	isVarName := e.tokenizer.Token().TokenType() == Identifier
	isUnary := e.tokenizer.Token().TokenType() == Symbol && (e.tokenizer.Token().Symbol() == '-' || e.tokenizer.Token().Symbol() == '~')
	isExpr := e.tokenizer.Token().TokenType() == Symbol && e.tokenizer.Token().Symbol() == '('

	// Check if we satisfy any of the above
	if !(isInt || isString || isKeyword || isVarName || isUnary || isExpr) {
		fmt.Println("Expected a term")
		return false
	}

	// Write opening
	e.write("<term>")
	e.indentVal++

	if isExpr {
		// We have an open parenthesis
		// Eat the symbol
		e.writeSymbol()

		// Now we should have an expression
		if !e.compileExpression() {
			fmt.Println("Expected expression")
			return false
		}

		// Now we should have a close parenthesis
		if !e.writeSpecSymbol(')') {
			fmt.Println("Expected close parenthesis")
			return false
		}

	} else if isUnary {
		// We have a unary operation
		// Eat the symbol
		e.writeSymbol()

		// Now we should have a term
		if !e.compileTerm() {
			fmt.Println("Expected term")
			return false
		}
	} else if isInt {
		// We have an integer constant
		e.write("<integerConstant> " + strconv.Itoa(e.tokenizer.Token().IntVal()) + " </integerConstant>")
		e.tokenizer.Advance()
	} else if isString {
		// We have a string constant
		e.write("<stringConstant> " + e.tokenizer.Token().StringVal() + " </stringConstant>")
		e.tokenizer.Advance()
	} else if isKeyword {
		// We have a keyword constant
		e.writeKeyword()
	} else if isVarName {
		// This could be a variable name, array, or subroutine call
		// We know we have an identifier so let's grab that string
		e.writeIdentifier()

		// Now we may have a period
		if e.writeSpecSymbol('.') {
			// We have a period so an identifier should follow
			if !e.writeIdentifier() {
				fmt.Println("Expected identifier")
				return false
			}
		}

		// Now we may have an open parenthesis or open bracket
		if e.writeSpecSymbol('[') {
			// In the open bracket we should have an expression
			if !e.compileExpression() {
				fmt.Println("Expected expression")
				return false
			}

			// Now we should have a closing bracket
			if !e.writeSpecSymbol(']') {
				fmt.Println("Expected closing bracket")
				return false
			}
		}

		// Now we may have an open parenthesis
		if e.writeSpecSymbol('(') {
			// We have an open parenthesis so we should have an expression list
			if !e.compileExpressionList() {
				fmt.Println("Expected expression list")
				return false
			}

			// Now we should have a closing parenthesis
			if !e.writeSpecSymbol(')') {
				fmt.Println("Expected closing parenthesis")
				return false
			}
		}

	}

	e.indentVal--
	e.write("</term>")

	return true
}

// compileExpressionList will write the xml for an expression list
func (e *Engine) compileExpressionList() bool {
	// Write opening
	e.write("<expressionList>")
	e.indentVal++

	// We have at least one expression if the next token isn't a closing paranthesis
	moreExpressions := e.tokenizer.Token().TokenType() != Symbol || e.tokenizer.Token().Symbol() != ')'
	for moreExpressions {
		// We should have an expression
		if !e.compileExpression() {
			fmt.Println("Expected expression")
			return false
		}

		// Now we may have a comma, which will indicate if we have more expressions
		moreExpressions = e.writeSpecSymbol(',')
	}

	e.indentVal--
	e.write("</expressionList>")

	return true

}

// WriteXML will process the jack file and write the results to an xml file matching the name
func (e *Engine) WriteXML() {
	// Create the xml file
	outputPath := strings.Replace(e.inputPath, ".jack", ".xml", 1)
	var err error
	e.outputFile, err = os.Create(outputPath)

	if err != nil {
		panic(err)
		return
	}

	// Close the file when we are done, panic if we hit an error
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(e.outputFile)

	// Now compile the class
	if !e.compileClass() {
		fmt.Println("Failed to compile... See above for details.")
		return
	}

}
