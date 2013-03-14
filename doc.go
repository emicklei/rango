/*
Rango is a REPL (Read-Evaluate-Print-Loop) tool in Go for Go.


Get the sources
		go get -v github.com/emicklei/rango

Install rango
		go install ...rango

Rango shell commands
		.q(uit)		exit rango
		.v(ars)		show all variable names
		.s(ource)	print the source entered since startup		
		.u			undo the last entry (e.g. to fix a compiler error)
		<name>		print a value when entered a known variable name

Features
	import declaration
	(almost) any go source that you can put inside the main() function

TODO
	function declarations
	multi variable declarations per entry
	multi import declarations per entry
	multi statements per entry

How it is made

Rango uses a generate-compile-run loop.
Successively, for each new command line entry, a new program is generated in Go, compiled in Go and run on your machine.
Any compiler error of the generated source is captured and printed by rango.
The output (stdout and stderr) of the generated program is captured and printed by rango.

(c) 2013, Ernest Micklei. MIT License
*/
package main
