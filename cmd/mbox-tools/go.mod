module github.com/aquilax/mbox-tools/cmd/mbox-tools

go 1.17

require (
	github.com/aquilax/mbox-tools/lib/mbox v0.0.0
	github.com/spf13/afero v1.8.2
	github.com/urfave/cli/v2 v2.4.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/aquilax/mbox-tools/lib/mbox => ../../lib/mbox
