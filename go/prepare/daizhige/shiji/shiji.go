package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func genParaID() string {
	return "para------"
}

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[1]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func splitSents(para string) []string {
	arr := []string{}
	split := RegSplit(para, "[。！：]”?")
	for i := range split {
		arr = append(arr, split[i])
	}
	return arr
}

func genSentID() string {
	return "sent------"
}

func main() {
	processShiji()
}

func testSplit() {
	para := "东至于海，登丸山，及岱宗。西至于空桐，登鸡头。南至于江，登熊、湘。"
	fmt.Println(para)
	split := splitSents(para)
	fmt.Println(split)
}

func processShiji() {
	dir := "D:/neoreads/data/test/"
	file, err := os.Open(dir + "test1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	out, err := os.Create(dir + "test1_out.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)

	// each line is a paragraph
	for scanner.Scan() {
		w.WriteString(genParaID())
		w.WriteString("|\n")
		para := scanner.Text()
		sents := splitSents(para)
		for _, sent := range sents {
			w.WriteString(genSentID())
			w.WriteString("|")
			w.WriteString(sent)
			w.WriteString("\n")
		}
		// paragraph delimiter
		w.WriteString("|\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
