package textseqchart

import (
	//"bufio"
	//"bytes"
	"fmt"
	"io"
	"os"
	//"strings"
	//"unicode"
)

func (invoke Invoke) arrow(index int) byte {
	var c byte
	if !invoke.IsSync && index%2 == 0 {
		c = ' '
	} else {
		c = '-'
	}
	return c
}

func BuildFile(out io.Writer, path string, ishtml bool) {
	seqdef, err := os.Open(path)
	defer seqdef.Close()
	//fmt.Printf("seqdef=%v, err=%v\n", seqdef, err)

	if err != nil {
		fmt.Println(os.Args[1], err)
		return
	}
	Build(out, seqdef, ishtml)
}

func BuildText(writer io.Writer, reader io.Reader) {
	Build(writer, reader, false)
}

func BuildHtml(writer io.Writer, reader io.Reader) {
	Build(writer, reader, true)
}

func Build(writer io.Writer, reader io.Reader, ishtml bool) {
	var (
		//title        string
		partnames    []string
		participants map[string]*Participant
		invacations  []*Invoke
	)
	_, partnames, participants, invacations = parse(reader)
	//fmt.Printf("\npartnames:%v \nparticipants:%v \ninvacations:%v\n", partnames, participants, invacations)

	first := participants[partnames[0]]
	first.Leftpos = 0
	first.Midpos = (len(first.Name)+3)/2 - 1 //两边有框；name6宽8取index3， name7宽9取index4

	//确定每个框的生命线位置
	for i := 0; i < len(partnames)-1; i++ {
		name1 := partnames[i]
		name2 := partnames[i+1]
		part1 := participants[name1]
		part2 := participants[name2]

		distance := len(name1)/2 + len(name2)/2 + 3 //两个框的最小距离
		for j := 0; j < len(invacations); j++ {
			invoke := invacations[j]
			msglen := distance
			//if ctrl,ok := invoke.(Control); ok{
			if invoke.Type != "" {
				if invoke.Invoker == part2	{
					//begin loop/opt/alt
					//msglen = len(invoke.Msg) + 2 + 5 //@note4
					msglen = len(invoke.Type) + 2 //@note4
				} //else if invoke.Target == part2 {	
					//end loop/opt/alt
				//}
			} else if invoke.Target.Seqnum > invoke.Invoker.Seqnum && invoke.Target == part2 {
				//箭头方向从左到右：partx(part1或跨part1左边的) ---> part2
				msglen = len(invoke.Msg) + 2
				msglen = msglen - (part1.Midpos - invoke.Invoker.Midpos)
			} else if invoke.Target.Seqnum < invoke.Invoker.Seqnum && invoke.Target == part1 {
				//箭头方向从右到左：part1 <--- partx(part2或跨part2右边的)
				msglen = len(invoke.Msg) + 2
				msglen = msglen - (invoke.Invoker.Midpos - part2.Midpos)
			} else if invoke.Target == invoke.Invoker && invoke.Target == part1 {
				//自调用：part1 ---> part1
				msglen = len(invoke.Msg) + 2 + 4 //note3 对应后面 note4
			}
			if msglen > distance {
				distance = msglen
			}
		}
		part2.Midpos = part1.Midpos + distance + 1
		part2.Leftpos = part2.Midpos - (len(name2)+3)/2 + 1
	}

	//fmt.Printf("title: %v \npartnames: %v \nparticipants:%v \ninvacations:%v \n", title, partnames, participants, invacations)

	chart := char2d{make([][]byte, 0)}
	//画箭头
	invokej := 4  //纵坐标
	for _, invoke := range invacations {
		//if ctrl,ok = invoke.(Control); ok{
		if invoke.Type != "" {
			invoke.Yops1 = invokej
			invokej += 1
			continue
		}

		invoker := participants[invoke.Invoker.Name]
		target := participants[invoke.Target.Name]
		invoke.Yops1 = invokej
		
		if target.Seqnum > invoker.Seqnum {
			//箭头方向从左到右：parta ---> partb
			writeMsg(&chart, invoke.Msg, invokej, invoker.Midpos+1)

			for x := invoker.Midpos + 1; x < target.Midpos-1; x++ {
				//b := byte(!invoke.IsSync && (x-invoker.Midpos)%2==0 ? ' ' : '-')
				c := invoke.arrow(x - invoker.Midpos)
				chart.Insert(invokej+1, x, c, true) //note2 和前面 note1 对应
			}
			chart.Insert(invokej+1, target.Midpos-1, '>', true)
			invokej += 2
		} else if target.Seqnum < invoker.Seqnum {
			//箭头方向从右到左：parta <--- partb
			writeMsg(&chart, invoke.Msg, invokej, target.Midpos+1)

			for x := invoker.Midpos - 1; x > target.Midpos; x-- {
				c := invoke.arrow(invoker.Midpos - x)
				chart.Insert(invokej+1, x, c, true) //note2 和前面 note1 对应
			}
			chart.Insert(invokej+1, target.Midpos+1, '<', true)
			invokej += 2
		} else { //target.Seqnum > invoker.Seqnum
			//自调用：parta ---> parta
			writeMsg(&chart, invoke.Msg, invokej+1, target.Midpos+1+4) //note4 对应前面 note3
			chart.Insert(invokej, target.Midpos+0, '-', true)          //|
			chart.Insert(invokej, target.Midpos+1, '-', true)          //+---.
			chart.Insert(invokej, target.Midpos+2, '-', true)          //|   |
			chart.Insert(invokej, target.Midpos+3, '-', true)          //|<--'
			chart.Insert(invokej, target.Midpos+4, '.', true)
			chart.Insert(invokej+1, target.Midpos+4, '|', true)
			chart.Insert(invokej+2, target.Midpos+4, '\'', true)
			chart.Insert(invokej+2, target.Midpos+3, '-', true)
			chart.Insert(invokej+2, target.Midpos+2, '-', true)
			chart.Insert(invokej+2, target.Midpos+1, '<', true)
			invokej += 3
		}
		invoke.Yops2 = invokej
	}
	
	//画Control（loop/opt/alt)
	for _, invoke := range invacations {
		if invoke.Type == "" || invoke.Invoker == nil {
			continue
		}
		
		len := len(invoke.Type)
		ctrlYops1 := invoke.Yops1
		ctrlYops2 := invoke.Peer.Yops1
		ctrlXops1 := invoke.Invoker.Midpos-len-1
		for x:= ctrlXops1; x <= invoke.Invoker.Midpos; x++ {
			chart.Insert(ctrlYops1, x, '-', true)
			chart.Insert(ctrlYops2, x, '-', true)
		}
		//chart.Insert(ctrlYops1, invoke.Invoker.Midpos-2, '<', true)
		//chart.Insert(ctrlYops2, invoke.Invoker.Midpos-2, '>', true)
		chart.Insert(ctrlYops1, ctrlXops1, '.', true)
		chart.Insert(ctrlYops2, ctrlXops1, '\'', true)
		writeMsg(&chart, invoke.Type, ctrlYops1+1, ctrlXops1+1)
		
		for y:= ctrlYops1+1; y < ctrlYops2; y++ {
			chart.Insert(y, ctrlXops1, '|', true)
		}

	}
	
	//maxj := len(invacations)*2 + 5 //note1: 和后面 note2 对应
	for name, part := range participants {
		x1 := part.Leftpos
		x2 := part.Leftpos + 1 + len(part.Name)
		//画方框
		printRectagle(0, x1, x2, name, &chart)
		printRectagle(invokej+1, x1, x2, name, &chart)
		//画生命线
		for j := 2; j <= invokej+1; j++ {
			chart.Insert(j, part.Midpos, '|', false)
		}
	}

	//for
	//fmt.Println(chart)
	if ishtml {
		writer.Write([]byte(chart.String3()))
		//fmt.Fprintf(writer, "%v", chart)
	} else {
		fmt.Fprintf(writer, "%v", chart)
	}
}

func writeMsg(chart *char2d, msg string, row int, offset int) {
	for x, b := range []byte(msg) {
		chart.Insert(row, offset+x, b, true)
	}
}

//画方框
func printRectagle(y, x1, x2 int, name string, chart *char2d) {
	for x := x1 + 1; x < x2; x++ {
		chart.Insert(y, x, '-', true)
		chart.Insert(y+2, x, '-', true)
	}
	//for x, b := range name {
	//	chart.Insert(y+1, x1+1+x, byte(b), true)
	//}
	writeMsg(chart, name, y+1, x1+1)
	chart.Insert(y+0, x1, '+', true)
	chart.Insert(y+0, x2, '+', true)
	chart.Insert(y+2, x1, '+', true)
	chart.Insert(y+2, x2, '+', true)
	chart.Insert(y+1, x1, '|', true)
	chart.Insert(y+1, x2, '|', true)
}
