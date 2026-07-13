/*
Package goroutines supports discovering the states of go routines.

This package encapsulates the gory details of parsing the Go runtime's textual
stack traces in order to extract the go routine details, such as ID,
thread-locked status, synctest bubble ID, and (when enabled using GODEBUG=) go
routine labels.

# Go Routine Stack Dump Header Line Syntax

Caveat Emptor: the following [ABNF] syntax (also [ABNF Update]) is inofficial
and can change with each minor Go release.

	goroutine-header =
	    "goroutine" SP go-id
	    [ runtime ]
	    SP "["
	    status
	    [ SP "(leaked)" ]
	    [ SP "(scan)" ]
	    [ SP "(durable)" ]
	    [ "," SP 1*DIGIT SP "minutes" ]
	    [ "," SP "locked to thread" ]
	    [ "," SP "synctest bubble" SP bubble-id ]
	    [ SP "labels:" labels ]
	    "]:" LF

	runtime =
	    SP "gp=" pointer
	    SP ( "m=nil" / ( "m=" machine-id SP "mp=" pointer ) )

	labels =
	    SP "{" label *( "," SP label ) "}"

	label =
	    quoted-string ":" SP quoted-string

	quoted-string = ...

	go-id = 1*DIGIT
	machine-id = 1*DIGIT
	bubble-id = 1*DIGIT
	pointer = "0x" 1*HEXDIGIT

	LF = %x0A
	SP = " "
	DIGIT = %x30-39
	HEXDIGIT = DIGIT / %x41-46 / %x61-66

[ABNF]: https://www.rfc-editor.org/info/rfc5234/
[ABNF Update]: https://www.rfc-editor.org/info/rfc7405/
*/
package goroutines
