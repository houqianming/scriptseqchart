package timeline

import (
	"fmt"
	"io"
	"time"
)

//1226
//0105

func BuildFile(out io.Writer, path string, ishtml bool) {
	fmt.Println("开始画图")
	p := fmt.Println
	t := time.Now()
	p(t.Format("3:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-02-02T15:04:05.999999-07:00"))
	p(t.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	the_time, err := time.Parse("2006-01-02", "2014-12-26")
	if err == nil {
		p(the_time.Format("2006-02-02 15:04:05"))
	}
}
