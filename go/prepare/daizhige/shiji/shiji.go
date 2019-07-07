package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/neoreads-backend/go/server/repositories"

	"github.com/neoreads-backend/go/prepare/models"

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

var db *sqlx.DB
var bookRepo *repositories.BookRepo

func genBookID() string {
	bookGen := util.NewN64Generator(8)
	bookid := bookGen.Next()
	genCount := 0
	for bookRepo.ContainsBookID(bookid) {
		bookid := bookGen.Next()
		genCount++
		if genCount > 100 {
			log.Fatalf("Error generating bookid:%s", bookid)
			break
		}
	}
	return bookid
}

func createBookDir(dir string, bookid string) string {
	bookDir := filepath.Join(dir, "books", bookid[0:4], bookid)
	os.MkdirAll(bookDir, os.ModePerm)
	return bookDir
}

func testAddBook() {
	dir := "D:/neoreads/data/"
	db, err := sqlx.Connect("postgres", "user=postgres dbname=neoreads sslmode=disable password=123456")
	if err != nil {
		log.Fatalf("init db failed: %s\n", err)
	}
	bookRepo = repositories.NewBookRepo(db, dir+"books/")

	toc := models.NewToc("ADXZIQNZ", "测试书籍")
	toc.AddChapter("001", "Chapter 1")
	toc.AddChapter("002", "Chapter 2")
	toc.AddChapter("003", "Chapter 3")

	bookRepo.AddBook(toc)
}
func processShiji() {
	util.InitSeed()
	dir := "D:/neoreads/data/"
	db, err := sqlx.Connect("postgres", "user=postgres dbname=neoreads sslmode=disable password=123456")
	if err != nil {
		log.Fatalf("init db failed: %s\n", err)
	}
	bookRepo = repositories.NewBookRepo(db, dir+"books/")
	bookid := genBookID()
	log.Printf("new bookid: %s\n", bookid)

	bookDir := createBookDir(dir, bookid)
	bookFile := "史记.txt"
	toc := processBook(dir, bookFile, bookDir, bookid)
	bookRepo.AddBook(toc)

	// for each book convert it into .md with ids
	for _, f := range util.FindFile(bookDir, "*.txt") {
		dir, txt := filepath.Split(f)
		name := util.StripExt(txt)
		processChapter(dir, txt, name+".md")
	}

}

func processBook(dir string, fname string, outdir string, bookid string) *models.Toc {
	chapGen := util.NewN64Generator(4)
	file, err := os.Open(filepath.Join(dir, "stage", fname))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	var title string
	var chapTitles []string
	var toc *models.Toc

	contentStart := false
	var curChap *models.Chapter
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
			toc = models.NewToc(bookid, title)
			curChap = models.NewChapter(filepath.Join(outdir, "book_toc.txt"))
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
				chapid := chapGen.Next()
				curChap = models.NewChapter(filepath.Join(outdir, chapid+".txt"))
				curChap.Title = title
				toc.AddChapter(chapid, title)
				chapCount++
			}
		} else {
			curChap.Add(line)
		}

		curChap.Save()
	}

	toc.SaveMD(filepath.Join(outdir, "__TOC.md"))

	log.Printf("title:%s\n", title)
	log.Printf("chap titles:%s\n", chapTitles)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return toc
}

func processChapter(dir string, infile string, outfile string) {
	chapGen := util.NewN64Generator(4)
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

	lineCount := 0
	// each line is a paragraph
	for scanner.Scan() {
		para := scanner.Text()
		if len(para) == 0 {
			continue
		}

		lineCount++
		// First line is title of the chapter
		if lineCount == 1 {
			w.WriteString("## ")
		}
		//log.Printf("para:{%s}\n", para)
		sents := splitSents(para)
		for _, sent := range sents {
			//log.Printf("{%s}\n", sent)
			w.WriteString(sent)
			w.WriteString("<sent id=\"")
			w.WriteString(chapGen.Next())
			w.WriteString("\"/>\n")
		}
		// paragraph id and delimter
		w.WriteString("<para id=\"")
		w.WriteString(chapGen.Next())
		w.WriteString("\"/>\n\n")
	}
	defer w.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//testSplit()
	//testGenN64()
	processShiji()
	//testAddBook()
}
