package textseqchart

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
	//"os"
	"strings"
	//"unicode"
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
	Yops1    int
	Yops2    int
	//for Control(loop/opt/alt)
	Type string
	Peer *Invoke
}

/*
type Control struct{
	Invoke
	Type string
	End *Control
}
*/

func (invoke Invoke) String() string {
	return fmt.Sprintf("%v -> %v : %v\n", invoke.Invoker, invoke.Target, invoke.Msg)
}

func isParticipantLine(line string) string {
	if strings.HasPrefix(line, "participant ") {
		return strings.TrimSpace(string(line[len("participant "):]))
	}
	return ""
}

func parse(reader io.Reader) (title string, partnames []string, participants map[string]*Participant, invacations []*Invoke) {
	buff := bufio.NewReader(reader)
	partnames = make([]string, 0)
	participants = make(map[string]*Participant)
	invacations = make([]*Invoke, 0)
	controlStack := make([]*Invoke, 2)
	seqnum := 0
	for {
		line, e := buff.ReadString('\n')
		if e != nil && io.EOF != e {
			break
		}
		//fmt.Printf("line:%v\n, e:%v\n", line, e)
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			if io.EOF == e {
				break
			} else {
				continue
			}
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

		if strings.HasPrefix(line, "begin "){
			//begin loop/alt/opt : do some thing
			ctrl_msg := strings.Split(line, ":")
			if len(ctrl_msg) < 2 {
				ctrl_msg = strings.Split(line, "：") //中文冒号
				if len(ctrl_msg) < 2 {
					fmt.Println("error format 3:" + line)
					return
				}
			}
			ctrlType := string(ctrl_msg[0][6:])
			beginControl := &Invoke{Type: ctrlType} //保留空格
			beginControl.Msg = ctrl_msg[1]
			controlStack = append(controlStack, beginControl)
			invacations = append(invacations, beginControl)
			continue
		}
		if strings.HasPrefix(line, "end "){
			//end loop/alt/opt
			ctrlType := string(line[4:])
			endControl := &Invoke{Type: ctrlType} //strings.TrimSpace
			beginControl := controlStack[len(controlStack)-1]
			controlStack = controlStack[0:len(controlStack)-1]
			if strings.TrimSpace(beginControl.Type) != strings.TrimSpace(endControl.Type) {
				fmt.Println(endControl.Type + " not match " + beginControl.Type)
				return
			}
			beginControl.Peer = endControl
			endControl.Peer = beginControl
			endControl.Target = beginControl.Invoker
			invacations = append(invacations, endControl)
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
			bm = strings.Split(abm[1], "：")
			if len(bm) < 2 {
				fmt.Println("error format 2	")
				return
			}
		}
		//fmt.Printf("%q \n", bm)

		invoke := &Invoke{IsSync: isSync}
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
		
		if len(controlStack) > 0 {
			//取loop/alt/opt范围之内最左边的Participant作为Control的Participant
			beginControl := controlStack[len(controlStack)-1]
			if beginControl == nil {
				continue
			}
			if beginControl.Invoker == nil || beginControl.Invoker.Seqnum > invoke.Invoker.Seqnum {
				beginControl.Invoker = invoke.Invoker
			}
			if beginControl.Invoker == nil || beginControl.Invoker.Seqnum > invoke.Target.Seqnum{
				beginControl.Invoker = invoke.Target
			}
			//TODO 嵌套
		}

		if io.EOF == e {
			break
		}
	}
	fmt.Printf("invacations: %v \n", invacations)
	return
}
