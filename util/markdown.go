package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Para struct {
	ID    string
	Sents []Sent
}

const (
	NORMAL_SENT = 0
	HEADER      = 1
	WHOLE_LINE  = 2
)

type Sent struct {
	ID      string
	Content string
	Type    int // 0: normal sentence, 1: header, 2: whole line
}

var sentEndPat = regexp.MustCompile("^.*\\[\\#....\\]$")
var paraEndPat = regexp.MustCompile("^.*\\{\\#....\\}$")

var listStartPat = regexp.MustCompile(`^\s*-\s.*`)
var olistStartPat = regexp.MustCompile(`^\s*[0-9]+\. `)

var codeBlockPat = regexp.MustCompile("^`{3}|~{3}")

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

	inCodeBlock := false

	var paras []Para
	curpara := &Para{}
	var tableHead string
	inTable := false
	inIndentBlock := false
	for i := range lines {
		//line := strings.TrimRightFunc(lines[i], unicode.IsSpace)
		line := lines[i]
		if paraEndPat.MatchString(line) {
			l := len(line)
			paraid := line[l-5 : l-1]
			curpara.ID = paraid
			line = line[0 : l-7]
			if strings.HasSuffix(line, " ") {
				line = line[0 : len(line)-1]
			}
			if strings.TrimSpace(line) == "" { // a seperate line para id (e.g. after table or code)
				continue
			}
		}

		if idx := codeBlockPat.FindStringIndex(line); idx != nil {
			inCodeBlock = !inCodeBlock
			if inCodeBlock { // entering code block
				sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
				curpara.Sents = append(curpara.Sents, sent)
			} else { // exiting code block
				sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
				curpara.Sents = append(curpara.Sents, sent)
				paras = endPara(curpara, paras, paraIDGen)
			}
			continue
		}
		if inCodeBlock {
			sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
			curpara.Sents = append(curpara.Sents, sent)
			continue
		}
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			if !inIndentBlock {
				inIndentBlock = true
			}
			sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
			curpara.Sents = append(curpara.Sents, sent)
			continue
		} else {
			if inIndentBlock {
				inIndentBlock = false
				paras = endPara(curpara, paras, paraIDGen)
			}
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { // empty line, a new para
			if inIndentBlock { // unless it is in an indent block
				sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
				curpara.Sents = append(curpara.Sents, sent)
			} else {
				if len(curpara.Sents) > 0 {
					paras = endPara(curpara, paras, paraIDGen)
				}
				paras = append(paras, *curpara)
				curpara = &Para{}
			}
		} else if listStartPat.MatchString(line) { // unordered list
			sents := ParseSents(line, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)
			paras = endPara(curpara, paras, paraIDGen)
		} else if idx := olistStartPat.FindStringIndex(line); idx != nil { // ordered list
			log.Printf("match olist: %v", idx)
			header := Sent{ID: "", Content: line[idx[0]:idx[1]]}
			curpara.Sents = append(curpara.Sents, header)
			content := line[idx[1]:len(line)]
			log.Printf("sub content:%v\n", content)
			sents := ParseSents(content, sentIDGen)
			curpara.Sents = append(curpara.Sents, sents...)
			log.Printf("sents:%#v\n", curpara.Sents)
			paras = endPara(curpara, paras, paraIDGen)
		} else if strings.Contains(line, "|") { // may be table
			if !inTable && tableHead == "" {
				tableHead = line
				inTable = true
				continue
			}
			if inTable {
				if tableHead != "" { // second line of table
					sent := Sent{ID: "", Content: tableHead, Type: WHOLE_LINE}
					curpara.Sents = append(curpara.Sents, sent)
					tableHead = ""
				}
				sent := Sent{ID: "", Content: line, Type: WHOLE_LINE}
				curpara.Sents = append(curpara.Sents, sent)
			}
		} else {
			if inTable { //  end of table
				inTable = false
			}
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
	for i := range paras {
		fmt.Printf("%#v\n", paras[i])
	}
	return Paras2Text(paras)
}

func Paras2Text(paras []Para) string {
	fmt.Printf("paras:%v\n", len(paras))
	sb := strings.Builder{}

	for i := range paras {
		p := paras[i]

		if i > 0 {
			sb.WriteString("\n")
		}

		sents := p.Sents
		for j := range sents {
			sent := sents[j]
			sb.WriteString(sent.Content)
			if sent.ID != "" {
				if strings.HasSuffix(sent.Content, " ") {
					sb.WriteString("[#" + sent.ID + "]")

				} else {
					sb.WriteString(" [#" + sent.ID + "]")
				}
			}
			if sent.Type == WHOLE_LINE {
				sb.WriteString("\n")
			}
		}
		// write paraid
		if len(p.Sents) > 0 && p.ID != "" {
			if p.Sents[len(p.Sents)-1].Type == WHOLE_LINE {
				sb.WriteString("{#" + p.ID + "}")
			} else {
				sb.WriteString(" {#" + p.ID + "}")
			}
		}
	}

	return sb.String()
}
