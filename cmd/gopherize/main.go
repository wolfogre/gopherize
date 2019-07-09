package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wolfogre/gopher/internal/pkg/gopherize"
)

var (
	n = flag.Int("n", 1, "number")
	v = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()

	gopherize.SetVerbose(*v)

	for i := 0; i < *n; i++ {
		content, filename, err := gopherize.RandomImage()
		if err != nil {
			panic(err)
		}

		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		if _, err := file.Write(content); err != nil {
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}
		fmt.Println(filename)
	}
}
