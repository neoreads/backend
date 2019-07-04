package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/neoreads-backend/go/util"
)

// RegSplit Split a paragraph into sentences.
// Note: unlike the standard regexp.Split, the delimeters are retained
func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	haslast := len(indexes) > 0 && indexes[len(indexes)-1][1] == len(text)
	log.Println(indexes)
	laststart := 0
	var result []string
	if haslast {
		result = make([]string, len(indexes))
	} else {
		result = make([]string, len(indexes)+1)
	}
	for i, element := range indexes {
		// retain the delimiter
		result[i] = text[laststart:element[1]]
		laststart = element[1]
	}
	if !haslast {
		result[len(indexes)] = text[laststart:len(text)]
	}
	return result
}

// splitSents Split sentences with Chinese punctuations
func splitSents(para string) []string {
	arr := []string{}
	split := RegSplit(para, "([。！：]”?)+")
	for i := range split {
		arr = append(arr, split[i])
	}
	return arr
}

func testGenN64() {
	util.InitSeed()
	ng := util.NewN64Generator(3)
	for i := 0; i < 10000; i++ {
		fmt.Println(ng.Next())
	}
}

func testSplit() {
	//para := "东至于海，登丸山，及岱宗。西至于空桐，登鸡头。南至于江，登熊、湘。！"
	para := "你好"
	fmt.Println(para)
	split := splitSents(para)
	for _, line := range split {
		fmt.Printf("{%s}\n", line)
	}
}

func processShiji() {
	util.InitSeed()
	dir := "D:/neoreads/data/test/"
	processBook(dir, "史记.txt")
	//processChapter(dir, "test1.txt", "test1_out.txt")
}

type Chapter struct {
	Title    string
	fileName string
	lines    []string
}

func NewChapter(file string) *Chapter {
	return &Chapter{fileName: file}
}

func (c *Chapter) Add(line string) {
	c.lines = append(c.lines, line)
}

func (c *Chapter) Save() {
	out, err := os.Create(c.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	w.WriteString(c.Title)
	w.WriteString("\n")
	for _, line := range c.lines {
		w.WriteString(line)
		w.WriteString("\n")
	}
	w.Flush()
}

func processBook(dir string, fname string) {

	chapGen := util.NewN64Generator(4)
	file, err := os.Open(dir + fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	var title string
	var chapTitles []string
	contentStart := false
	var curChap *Chapter
	chapCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		lineCount++
		// first line is book title
		if lineCount == 1 {
			title = line
			curChap = NewChapter(dir + "toc.txt")
			curChap.Title = line
			chapCount++
			continue
		}
		// chapter titles
		if strings.HasPrefix(line, "卷") {
			title := line
			// toc
			if !contentStart {
				// title of first chapter comes again, meaning main content
				if len(chapTitles) > 0 && title == chapTitles[0] {
					log.Printf("Content start here:%s\n", title)
					contentStart = true
				}
			}
			if !contentStart {
				chapTitles = append(chapTitles, title)
				curChap.Add(title)
			} else { // start of each chapter
				curChap.Save()
				curChap = NewChapter(dir + "chapter_" + strconv.Itoa(chapCount) + ".txt")
				curChap.Title = title
				chapCount++
			}
		} else {
			curChap.Add(line)
		}

		curChap.Save()
	}

	log.Printf("title:%s\n", title)
	log.Printf("chap titles:%s\n", chapTitles)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	chapGen.Next()
}

func processChapter(dir string, infile string, outfile string) {
	paraGen := util.NewN64Generator(4)
	sentGen := util.NewN64Generator(4)
	file, err := os.Open(dir + infile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	out, err := os.Create(dir + outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)

	// each line is a paragraph
	for scanner.Scan() {
		para := scanner.Text()
		if len(para) == 0 {
			continue
		}
		w.WriteString(paraGen.Next())
		w.WriteString("|\n")
		log.Printf("para:{%s}\n", para)
		sents := splitSents(para)
		for _, sent := range sents {
			log.Printf("{%s}\n", sent)
			w.WriteString(sentGen.Next())
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

func main() {
	//testSplit()
	//testGenN64()
	processShiji()
}
