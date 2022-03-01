package main

import (
	"log"
	"os"

	"github.com/aquilax/mbox-tools/cmd"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	err := cmd.GetApp(fs).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
