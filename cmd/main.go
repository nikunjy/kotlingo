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
	imports := processor.GetImports()
	fmt.Println(imports)
	fmt.Println(processor.GetPackageName())
}
