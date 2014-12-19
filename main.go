package main

import (
	"fmt"
	"github.com/houqianming/scriptseqchart/textseqchart"
	"os"
)

func main() {
	fmt.Printf("args: %v \n", os.Args)
	if len(os.Args) < 2 {
		return
	}

	textseqchart.Build(os.Args[1])
}
