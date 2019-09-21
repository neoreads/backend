package main

import (
	"fmt"
	"strings"

	"github.com/neoreads/backend/util"
)

type Para struct {
	ID    string
	Sents []Sent
}

type Sent struct {
	ID      string
	Content string
}

func parseSents(line string, gen *util.N64Generator) []Sent {
	var sents []Sent
	ss := util.SplitSents(line)
	for i := range ss {
		sents = append(sents, Sent{ID: gen.Next(), Content: ss[i]})
	}
	return sents
}

func parsemd(s string) {

	paraIDGen := util.NewN64Generator(4)
	sentIDGen := util.NewN64Generator(4)

	lines := strings.Split(s, "\n")

	var paras []Para
	curpara := Para{paraIDGen.Next(), []Sent{}}
	for i := range lines {
		//line := strings.TrimRightFunc(lines[i], unicode.IsSpace)
		line := lines[i]
		trimmed := strings.TrimSpace(lines[i])
		fmt.Printf("%v: %v\n", i, line)
		if trimmed == "" { // empty line, a new para
			if len(curpara.Sents) > 0 {
				paras = append(paras, curpara)
				curpara = Para{paraIDGen.Next(), []Sent{}}
			}
		} else if strings.HasPrefix(line, "-") {
			sents := parseSents(line, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)
			paras = append(paras, curpara)
			curpara = Para{paraIDGen.Next(), []Sent{}}
		} else {
			sents := parseSents(line, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)
		}
	}

	fmt.Printf("%v, %#v\n", len(paras), paras)
}

func main() {
	fmt.Printf("Hello \n")
	s := `#Chapter Title

This Chapter covers:   

- First section   
- Second section   

And there is an end. I want to know. If you think about it: it will be ok.
This is the real end.

Where is the money? I haven't seen any one.
`
	parsemd(s)
}
