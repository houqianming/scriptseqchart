package main

import (
	"fmt"
	//"os/exec"
	"testing"
	"unicode"
)

func TestCharset(t *testing.T) {
	str := "a中文字b"
	fmt.Printf("str:%v, len:%v\n", str, len(str))

	bytes1 := str[:]
	fmt.Printf("bytes1:%q, len(bytes1):%v \n", bytes1, len(bytes1))

	bytes2 := make([]byte, 0)
	for _, b := range str {
		bytes2 = append(bytes2, byte(b))
	}
	fmt.Printf("bytes2:%q, len(bytes2):%v \n", bytes2, len(bytes2))

	bytes3 := make([]byte, 0)
	for _, b := range str[:] {
		bytes3 = append(bytes3, byte(b)) //rune
	}
	fmt.Printf("bytes3:%q, len(bytes3):%v \n", bytes3, len(bytes3))

	bytes4 := make([]byte, 0)
	temp := str[:]
	for i := 0; i < len(temp); i++ {
		bytes4 = append(bytes4, temp[i]) //rune
	}
	fmt.Printf("bytes4:%q, len(bytes4):%v \n", bytes4, len(bytes4))

	bytes5 := make([]byte, 0)
	for _, b := range []byte(str) {
		bytes5 = append(bytes5, b)
	}
	fmt.Printf("bytes5:%q, len(bytes5):%v \n", bytes5, len(bytes5))
}

func Viewlen(utf8str string) int {
	c := 0
	for _, b := range utf8str {
		if unicode.Is(unicode.Scripts["Han"], b) {
			c += 2
		} else {
			c += 1
		}
	}
	return c
}


func TestViewlen(t *testing.T) {
    fmt.Printf("viewlen:%v \n", Viewlen("a中文字c"));
}

/*
func TestCmd(t *testing.T) {
	command := exec.Command("cmd", "/c", "dir c:\\")
	fmt.Printf("%v\n", command)
	bytes, _ := command.Output()
	fmt.Printf("%q\n", bytes)
}
*/

func main2() {
	bytes := []byte{97, 41}
	fmt.Printf("%s", bytes[0])
	fmt.Printf("%s", bytes[1])
	fmt.Printf("%q\n", bytes)
}
