package main

import (
	//"github.com/houqianming/scriptseqchart/server"
	"github.com/houqianming/scriptseqchart/textseqchart"
	"os"
)

func main() {
	textseqchart.BuildFile(os.Stdout, os.Args[1], true)
}
