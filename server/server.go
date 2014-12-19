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
		w.Write([]byte("<div style={}>"))
		textseqchart.Build(w, strings.NewReader(script[0]))
		w.Write([]byte("</div>"))
	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "script:%v\n", script)
		fmt.Fprintf(w, "<form method=post action=/> <textarea name=script> </textarea><input type=submit value=commit> </form>", script)

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
