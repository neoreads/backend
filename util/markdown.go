package util

import (
	"fmt"
	"regexp"
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

var sentEndPat = regexp.MustCompile("^.*\\[\\#....\\]$")
var paraEndPat = regexp.MustCompile("^.*\\{\\#....\\}$")

var listStartPat = regexp.MustCompile("^\\w*-.*")

func ParseSents(line string, gen *N64Generator) []Sent {
	var sents []Sent
	ss := SplitSents(line)
	for i := range ss {
		sent := ss[i]
		if sentEndPat.MatchString(sent) {
			l := len(sent)
			sentid := sent[l-5 : l-1]
			sents = append(sents, Sent{ID: sentid, Content: sent[:l-7]})
		} else {
			sents = append(sents, Sent{ID: gen.Next(), Content: ss[i]})
		}
	}
	return sents
}

func endPara(para *Para, paras []Para, paraIDGen *N64Generator) []Para {

	if len(para.Sents) > 0 && para.ID == "" {
		para.ID = paraIDGen.Next()
	}
	paras = append(paras, *para)
	*para = Para{}
	return paras
}

func ParseMD(s string) []Para {

	paraIDGen := NewN64Generator(4)
	sentIDGen := NewN64Generator(4)

	lines := strings.Split(s, "\n")

	var paras []Para
	curpara := &Para{}
	for i := range lines {
		//line := strings.TrimRightFunc(lines[i], unicode.IsSpace)
		line := lines[i]
		if paraEndPat.MatchString(line) {
			l := len(line)
			paraid := line[l-5 : l-1]
			curpara.ID = paraid
			line = line[0 : l-8]
		}
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" { // empty line, a new para
			if len(curpara.Sents) > 0 {
				paras = endPara(curpara, paras, paraIDGen)
			}
			paras = append(paras, *curpara)
			curpara = &Para{}
		} else if listStartPat.MatchString(line) {
			sents := ParseSents(line, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)

			paras = endPara(curpara, paras, paraIDGen)
		} else {
			sents := ParseSents(line, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)
		}
	}

	if len(curpara.Sents) > 0 {
		paras = endPara(curpara, paras, paraIDGen)
	}

	return paras
}

func ApplyIDs(s string) string {
	paras := ParseMD(s)
	/*
		for i := range paras {
			fmt.Printf("%#v\n", paras[i])
		}
	*/
	return Paras2Text(paras)
}

func Paras2Text(paras []Para) string {
	fmt.Printf("paras:%v\n", len(paras))
	sb := strings.Builder{}

	for i := range paras {
		p := paras[i]

		sents := p.Sents
		for j := range sents {
			sent := sents[j]
			sb.WriteString(sent.Content)
			if strings.HasSuffix(sent.Content, " ") {
				sb.WriteString("[#" + sent.ID + "]")

			} else {
				sb.WriteString(" [#" + sent.ID + "]")
			}
		}
		if p.ID != "" {
			sb.WriteString(" {#" + p.ID + "}\n")
		} else {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
