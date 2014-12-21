package server

import (
	"fmt"
	"github.com/houqianming/scriptseqchart/textseqchart"
	"log"
	"net/http"
	"strings"
)

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	script := r.Form["script"]
	if len(script) > 0 {
		//w.Write([]byte(script[0]))
		fmt.Println(script[0])
		w.Write([]byte("<style>pre {line-height:54%;font-family: '宋体',monospace,consolas,'Courier New',fixed-width;}</style><pre>"))
		textseqchart.Build(w, strings.NewReader(script[0]), true)
		w.Write([]byte("</pre>"))
	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "请输入脚本（->表示方法调用）<br>\n")
		fmt.Fprintf(w, "<form method=post action=/> <textarea name=script cols=120 rows=20> </textarea><br><input type=submit value=commit> </form>")

	}
	//fmt.Fprint(w, "Hello!")
}

func Start() {
	var h Hello
	err := http.ListenAndServe("localhost:4000", h)
	if err != nil {
		log.Fatal(err)
	}
}
