package util

import (
	"fmt"
	"strings"
)

type Para struct {
	ID    string
	Sents []Sent
}

type Sent struct {
	ID      string
	Content string
}

func parseSents(line string, gen *N64Generator) []Sent {
	var sents []Sent
	ss := SplitSents(line)
	for i := range ss {
		sents = append(sents, Sent{ID: gen.Next(), Content: ss[i]})
	}
	return sents
}

func ParseMD(s string) []Para {

	paraIDGen := NewN64Generator(4)
	sentIDGen := NewN64Generator(4)

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

	return paras
}
