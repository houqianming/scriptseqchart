package main

import (
	"fmt"
	"github.com/houqianming/scriptseqchart/server"
	"github.com/houqianming/scriptseqchart/textseqchart"
	"os"
)

func main() {
	fmt.Printf("args: %v \n", os.Args)
	if len(os.Args) < 2 {
		return
	}

	switch len(os.Args) {
	case 2:
		switch os.Args[1] {
		case "-server":
			server.Start()
		default:
			textseqchart.BuildFile(os.Stdout, os.Args[1])
		}
	default:
		fmt.Print("error")

	}

}
