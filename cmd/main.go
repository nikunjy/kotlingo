package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikunjy/kotlingo"
)

func main() {
	processor, err := kotlingo.NewProcessor(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	imports, err := processor.GetImports()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imports)
}
