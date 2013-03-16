/*
Rango is a REPL (Read-Evaluate-Print-Loop) tool in Go for Go.


Get the sources
		go get -v github.com/emicklei/rango

Install rango
		go install ...rango

Run
		rango [projectname]

Rango shell commands
		.q(uit)		exit rango
		.v(ars)		show all variable names
		.s(ource)	print the source entered since startup		
		.u(undo)		the last entry (e.g. to fix a compiler error)
		<name>		print a value when entered a known variable name

Features
	import declaration
	(almost) any go source that you can put inside the main() function
	all entries are logged in a <projectname>.changes file. If such a file exists then rango will process its contents first.

Requirements
	Installation of Go 1+ SDK
	Because it depends on sh (e.g. bash) it only runs on a Go supported *nix OS

How it is made

Rango uses a generate-compile-run loop.
Successively, for each new command line entry, a new program is generated in Go, compiled in Go and run on your machine.
Any compiler error of the generated source is captured and printed by rango.
The output (stdout and stderr) of the generated program is captured and printed by rango.

TODO
	function declarations
	multi variable declarations per entry
	multi import declarations per entry
	multi statements per entry
	interpret compiler errors

(c) 2013, Ernest Micklei. MIT License
*/
package main
