package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Participant struct {
	Name    string
	Leftpos int
	Midpos  int
	Seqnum  int
}

func (part Participant) String() string {
	return fmt.Sprintf("%v,%v,%v\n", part.Name, part.Leftpos, part.Midpos)
}

type Invoke struct {
	Invoker *Participant
	Target  *Participant
	Msg     string
	IsSync  bool
	Yops    int
}

func (invoke Invoke) String() string {
	if invoke.Invoker != nil && invoke.Target != nil {
		return fmt.Sprintf("%v -> %v : %v\n", invoke.Invoker.Name, invoke.Target.Name, invoke.Msg)
	}
	return fmt.Sprint("no Invoker or Target\n")
}
func (invoke Invoke) arrow(index int) byte {
	var c byte
	if !invoke.IsSync && index%2 == 0 {
		c = ' '
	} else {
		c = '-'
	}
	return c
}

func isParticipantLine(line string) string {
	if strings.HasPrefix(line, "participant ") {
		return strings.TrimSpace(string(line[len("participant "):]))
	}
	return ""
}

func parse(file string) (title string, partnames []string, participants map[string]*Participant, invacations []Invoke) {
	seqdef, err := os.Open(file)
	defer seqdef.Close()
	//fmt.Printf("seqdef=%v, err=%v\n", seqdef, err)

	if err != nil {
		fmt.Println(os.Args[1], err)
		return
	}

	partnames = make([]string, 0)
	participants = make(map[string]*Participant)
	invacations = make([]Invoke, 0)
	buff := bufio.NewReader(seqdef)
	seqnum := 0
	for {
		line, e := buff.ReadString('\n')
		if e != nil || io.EOF == err {
			break
		}
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			continue
		}

		if name := isParticipantLine(line); name != "" {
			seqnum++
			part := Participant{Name: name, Seqnum: seqnum}
			partnames = append(partnames, name)
			participants[name] = &part
			continue
		}

		if strings.HasPrefix(line, "title ") {
			title = strings.TrimSpace(line[5:])
			continue
		}

		var isSync bool
		abm := strings.Split(line, "-->")
		if len(abm) == 2 {
			isSync = false
		} else {
			abm = strings.Split(line, "->")
			if len(abm) == 2 {
				isSync = true
			} else {
				fmt.Printf("error format 1: %v", line)
				return
			}
		}
		bm := strings.Split(abm[1], ":")
		if len(bm) < 2 {
			fmt.Println("error format 2	")
			return
		}
		//fmt.Printf("%q \n", bm)

		invoke := Invoke{IsSync: isSync}
		{
			name := strings.TrimSpace(abm[0])
			part, ok := participants[name]
			if !ok {
				seqnum++
				part = &Participant{Name: name, Seqnum: seqnum}
				partnames = append(partnames, name)
				participants[name] = part
			}
			invoke.Invoker = part
		}

		{
			name := strings.TrimSpace(bm[0])
			part, ok := participants[name]
			if !ok {
				seqnum++
				part = &Participant{Name: name, Seqnum: seqnum}
				partnames = append(partnames, name)
				participants[name] = part
			}
			invoke.Target = part
		}

		invoke.Msg = strings.Join(bm[1:], ":")
		invacations = append(invacations, invoke)

	}
	return
}

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
	//if old != ' ' {
	//	fmt.Printf("old=%q, b=%q\n", old, b)
	//}
	if override || old == ' ' {
		if (old == '-' && b == '|') || (old == '|' && b == '-') {
			chart.bytes[i][j] = '+'
			//fmt.Printf("old=%q, b=%q, res= %q \n", old, b, chart.bytes[i][j] )
		} else {
			chart.bytes[i][j] = b
		}
	}
}

func (chart char2d) String() string {
	for i := 0; i < len(chart.bytes); i++ {

		for j := 0; j < len(chart.bytes[i]); j++ {
			fmt.Printf("%s", string(chart.bytes[i][j]))
		}
		fmt.Print("\n")

		//fmt.Printf("%q\n", chart.bytes[i])
	}
	return ""
}

func main() {
	//argnum := os.Args
	fmt.Printf("args: %v \n", os.Args)
	if len(os.Args) < 2 {
		return
	}

	var (
		//title        string
		partnames    []string
		participants map[string]*Participant
		invacations  []Invoke
	)
	_, partnames, participants, invacations = parse(os.Args[1])

	first := participants[partnames[0]]
	first.Leftpos = 0
	first.Midpos = (len(first.Name)+3)/2 - 1 //两边有框；name6宽8取index3， name7宽9取index4

	//确定位置
	for i := 0; i < len(partnames)-1; i++ {
		name1 := partnames[i]
		name2 := partnames[i+1]
		part1 := participants[name1]
		part2 := participants[name2]

		distance := len(name1)/2 + len(name2)/2 + 2 //两个框的最小距离
		for j := 0; j < len(invacations); j++ {
			if invacations[j].Target.Seqnum > invacations[j].Invoker.Seqnum && invacations[j].Target == part2 {
				msglen := len(invacations[j].Msg) + 2
				msglen = msglen - (part1.Midpos - invacations[j].Invoker.Midpos)
				if msglen > distance {
					distance = msglen
				}
			}
			if invacations[j].Target.Seqnum < invacations[j].Invoker.Seqnum && invacations[j].Target == part1 {
				msglen := len(invacations[j].Msg) + 2
				msglen = msglen - (invacations[j].Invoker.Midpos - part2.Midpos)
				if msglen > distance {
					distance = msglen
				}
			}
			if invacations[j].Target == invacations[j].Invoker && invacations[j].Target == part1 {
				msglen := len(invacations[j].Msg) + 2 + 4 //note3 对应后面 note4
				if msglen > distance {
					distance = msglen
				}
			}

		}
		part2.Midpos = part1.Midpos + distance + 1
		part2.Leftpos = part2.Midpos - (len(name2)+3)/2 + 1
	}

	//fmt.Printf("title: %v \npartnames: %v \nparticipants:%v \ninvacations:%v \n", title, partnames, participants, invacations)

	chart := char2d{make([][]byte, 0)}
	//画箭头
	invokej := 4
	for _, invoke := range invacations {
		invoker := participants[invoke.Invoker.Name]
		target := participants[invoke.Target.Name]

		if target.Seqnum > invoker.Seqnum {
			for x, b := range invoke.Msg {
				chart.Insert(invokej, invoker.Midpos+1+x, byte(b), true)
			}

			for x := invoker.Midpos + 1; x < target.Midpos-1; x++ {
				//b := byte(!invoke.IsSync && (x-invoker.Midpos)%2==0 ? ' ' : '-')
				c := invoke.arrow(x - invoker.Midpos)
				chart.Insert(invokej+1, x, c, true) //note2 和前面 note1 对应
			}
			chart.Insert(invokej+1, target.Midpos-1, '>', true)
			invokej += 2
		} else if target.Seqnum < invoker.Seqnum {
			for x, b := range invoke.Msg {
				chart.Insert(invokej, target.Midpos+1+x, byte(b), true)
			}

			for x := invoker.Midpos - 1; x > target.Midpos; x-- {
				c := invoke.arrow(invoker.Midpos - x)
				chart.Insert(invokej+1, x, c, true) //note2 和前面 note1 对应
			}
			chart.Insert(invokej+1, target.Midpos+1, '<', true)
			invokej += 2
		} else { //target.Seqnum > invoker.Seqnu
			for x, b := range invoke.Msg {
				chart.Insert(invokej+1, target.Midpos+1+4+x, byte(b), true) //note4 对应前面 note3
			}
			chart.Insert(invokej, target.Midpos+0, '-', true) //|
			chart.Insert(invokej, target.Midpos+1, '-', true) //+---.
			chart.Insert(invokej, target.Midpos+2, '-', true) //|   |
			chart.Insert(invokej, target.Midpos+3, '-', true) //|<--'
			chart.Insert(invokej, target.Midpos+4, '.', true)
			chart.Insert(invokej+1, target.Midpos+4, '|', true)
			chart.Insert(invokej+2, target.Midpos+4, '\'', true)
			chart.Insert(invokej+2, target.Midpos+3, '-', true)
			chart.Insert(invokej+2, target.Midpos+2, '-', true)
			chart.Insert(invokej+2, target.Midpos+1, '<', true)
			invokej += 3
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
			chart.Insert(j, part.Midpos, '|', true)
		}
	}

	//for
	fmt.Println(chart)
}

//画方框
func printRectagle(y, x1, x2 int, name string, chart *char2d) {
	for x := x1 + 1; x < x2; x++ {
		chart.Insert(y, x, '-', true)
		chart.Insert(y+2, x, '-', true)
	}
	for x, b := range name {
		chart.Insert(y+1, x1+1+x, byte(b), true)
	}
	chart.Insert(y+0, x1, '+', true)
	chart.Insert(y+0, x2, '+', true)
	chart.Insert(y+2, x1, '+', true)
	chart.Insert(y+2, x2, '+', true)
	chart.Insert(y+1, x1, '|', true)
	chart.Insert(y+1, x2, '|', true)
}
