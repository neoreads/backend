package models

import (
	"bufio"
	"log"
	"os"
)

type Toc struct {
	BookID string
	Title  string
	Items  []TocItem
}

type TocItem struct {
	ChapID string
	Title  string
}

func NewToc(bookid string, title string) *Toc {
	return &Toc{
		BookID: bookid,
		Title:  title,
	}
}

func (t *Toc) AddChapter(chapid string, title string) {
	t.Items = append(t.Items, TocItem{ChapID: chapid, Title: title})
}

// SaveMD save toc to a .md file
func (t *Toc) SaveMD(path string) {
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	// Book Title
	w.WriteString("# ")
	w.WriteString(t.Title)
	w.WriteString("\n\n")

	// List of toc items
	for _, item := range t.Items {
		w.WriteString("- [")
		w.WriteString(item.Title)
		w.WriteString("](")
		w.WriteString(item.ChapID)
		w.WriteString(".md)\n")
	}
}

type Book struct {
	ID    string
	Title string
}

func NewBook(title string) *Book {
	return &Book{Title: title}
}

type Chapter struct {
	ID       string
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
