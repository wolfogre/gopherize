package main

import (
	"fmt"
	"os"

	"github.com/wolfogre/gopher/internal/app/gopherize"
)

func main() {
	content, filename, err := gopherize.RandomImage()
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err, _ := file.Write(content); err != nil {
		panic(err)
	}
	fmt.Println(filename)
}
