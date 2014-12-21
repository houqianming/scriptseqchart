package textseqchart

import (
	//"bufio"
	"bytes"
	//"fmt"
	//"io"
	//"os"
	//"strings"
	"unicode"
)

type char2d struct {
	bytes [][]byte
}

func (chart *char2d) Insert(i, j int, b byte, override bool) {
	leni := len(chart.bytes)
	if i >= leni {
		for k := leni; k <= i; k++ {
			chart.bytes = append(chart.bytes, make([]byte, 0))
		}
	}
	lenj := len(chart.bytes[i])
	if j >= lenj {
		for k := lenj; k <= j; k++ {
			chart.bytes[i] = append(chart.bytes[i], ' ')
		}
	}
	old := chart.bytes[i][j]
	if (old == '-' && b == '|') || (old == '|' && b == '-') {
		chart.bytes[i][j] = '+'
		//fmt.Printf("old=%q, b=%q, res= %q \n", old, b, chart.bytes[i][j] )
	} else if override || old == ' ' {
		chart.bytes[i][j] = b
	}
}

func (chart char2d) String1() string {
	var buffer bytes.Buffer
	for i := 0; i < len(chart.bytes); i++ {
		buffer.Write(chart.bytes[i])
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
func (chart char2d) String2() string {
	var buffer bytes.Buffer
	for i := 0; i < len(chart.bytes); i++ {
		for _, c := range string(chart.bytes[i]) {
			buffer.WriteRune(c)
			if unicode.Is(unicode.Scripts["Han"], c) {
				//因汉字go内部3个字节，文本显示时等宽占2个字符，所以多些一个空格
				buffer.WriteByte(' ')
				//fmt.Printf("c:%v\n", c)
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
func (chart char2d) String3() string {
	var buffer bytes.Buffer
	for i := 0; i < len(chart.bytes); i++ {
		for _, c := range string(chart.bytes[i]) {
			switch c {
			case ' ':
				buffer.WriteString("&nbsp;")
			case '<':
				buffer.WriteString("&lt;")
			case '>':
				buffer.WriteString("&gt;")
			default:
				buffer.WriteRune(c)
			}
			if unicode.Is(unicode.Scripts["Han"], c) {
				//因汉字go内部3个字节，文本显示时等宽占2个字符，所以多写一个空格
				buffer.WriteString("&nbsp;")
				//fmt.Printf("c:%v\n", c)
			}
		}
		buffer.WriteString("<br>\n")
	}
	return buffer.String()
}
func (chart char2d) String() string {
	return chart.String2()
}
