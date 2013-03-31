/*
Rango is a REPL (Read-Evaluate-Print-Loop) tool in Go for Go.


Get the sources
		go get -v github.com/emicklei/rango

Install
		go install ...rango

Run
		rango [projectname]

Example session
	> rango
	[rango] .q = quit, .v = variables, .s = source, .u = undo
	> m,y := "rango the chameleon", 2012
	rango the chameleon,2012
	> import "strings"
	> m = strings.ToUpper(m)
	RANGO THE CHAMELEON
	> y+1
	2013

Commands
		.q(uit)		exit rango
		.v(ars)		show all variable names
		.s(ource)	print the source entered since startup		
		.u(undo)	the last entry
		<expression>		print a value when entered an expression

Features
	import declaration
	(almost) any go source that you can put inside the main() function
	if <projectname> is given on startup then
		if a <projectname>.changes file exists then rango will process its contents first.
		all entries are logged in a <projectname>.changes file.

Requirements
	Installation of Go 1+ SDK
	Because it depends on sh (e.g. bash) it only runs on a Go supported *nix OS


How it is made

Rango uses a generate-compile-run loop.
Successively, for each new command line entry, a new program is generated in Go, compiled in Go and run on your machine.
Any compiler error of the generated source is captured and printed by rango.
The output (stdout and stderr) of the generated program is captured and printed by rango.

Todo

	interpret compiler errors and translate line numbers
	use goreadline? termbox-go? for better cursor handling (up,down,complete...)

(c) 2013, Ernest Micklei. MIT License
*/
package main
