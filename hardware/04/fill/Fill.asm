// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

(LOOP)
	// Look for key press
	@KBD
	D=M

	@FILL
	D;JNE

	// Look for key press since D is maintained
	@KBD
	D=M

	@CLEAR
	D;JEQ

	@LOOP
	0;JMP


(FILL)
	// Initialize screen count
	@scrnPtr
	M=0

	@FLOOP
	0;JMP

(FLOOP) // Fill loop
	@scrnPtr
	D=M
	@SCREEN
	D=D+A
	@currScreen
	A=D
	M=-1

	// Increment loop then check
	@scrnPtr
	M=M+1
	D=M
	@8192
	D=D-A

	// Jump to head if loop is done
	@LOOP
	D;JEQ

	// Continue looping otherwise
	@FLOOP
	D;JNE

(CLEAR)
	// Initialize screen count
	@scrnPtr
	M=0

	@CLOOP
	0;JMP

(CLOOP) // Clear loop
	@scrnPtr
	D=M
	@SCREEN
	D=D+A
	@currScreen
	A=D
	M=0

	// Increment loop then check
	@scrnPtr
	M=M+1
	D=M
	@8192
	D=D-A

	// Jump to head if loop is done
	@LOOP
	D;JEQ

	// Continue looping otherwise
	@CLOOP
	D;JNE
