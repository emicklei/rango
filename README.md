rango
=====

a REPL program for the Go language

fetch
---
go get -v github.com/emicklei/rango

build
---
<code>go build rango</code>

run
---
./rango


commands
---
Inside the rango shell you can

* .q(uit)	exit rango
* .v(ars) 	show all variable names
* .u(ndo)	undo the last entry (can be repeated)
* print a value when entered a known variable name

features
----
* import declaration
* (almost) any go source you can put inside the main() function

todo
---
Currently rango can not handle

* function declarations
* multi variable declarations
* multi import declarations

how it is made
---
Rango uses a generate-compile-gorun loop. Successively, for each new command line entry, a new program is generated in Go, compiled in Go and run on your machine. The output (stdout and stderr) of the generated program is captured and printed by rango.