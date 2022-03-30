package main

import (
	"log"
	"os"

	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	err := getApp(fs).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
