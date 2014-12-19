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
	Yops    int
}

func (invoke Invoke) String() string {
	if invoke.Invoker != nil && invoke.Target != nil {
		return fmt.Sprintf("%v -> %v : %v\n", invoke.Invoker.Name, invoke.Target.Name, invoke.Msg)
	}
	return fmt.Sprint("no Invoker or Target\n")
}

func isParticipantLine(line string) string {
	if strings.HasPrefix(line, "participant ") {
		return strings.TrimSpace(string(line[len("participant "):]))
	}
	return ""
}

func parse(reader io.Reader) (title string, partnames []string, participants map[string]*Participant, invacations []Invoke) {
	buff := bufio.NewReader(reader)
	partnames = make([]string, 0)
	participants = make(map[string]*Participant)
	invacations = make([]Invoke, 0)
	seqnum := 0
	for {
		line, e := buff.ReadString('\n')
		if e != nil && io.EOF != e {
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

		if io.EOF == e {
			break
		}
	}
	return
}
