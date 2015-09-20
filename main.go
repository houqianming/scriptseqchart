package main

import (
	"fmt"
	"github.com/houqianming/scriptseqchart/html5/timeline"
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
		case "-timeline":
			timeline.BuildFile(os.Stdout, os.Args[1], false)
		default:
			textseqchart.BuildFile(os.Stdout, os.Args[1], false)
		}
	default:
		fmt.Print("error")

	}

}
