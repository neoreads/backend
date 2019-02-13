package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// 判断正文开始的分隔行
func isStartDelim(line string) bool {
	return strings.HasPrefix(line, "*** START OF")
}

// 判断正文结尾的分隔行
func isEndDelim(line string) bool {
	return strings.HasPrefix(line, "*** END OF")
}

// 判断是否为章节标题行
func isChapterTitle(line string) bool {
	return strings.HasPrefix(line, "Chapter")
}

// 提取章节编号
func extractChapterTitle(line string) string {
	return strings.TrimSpace(line[len("Chapter"):])
}

// 用来收集一个章节的数据结构
type chapter struct {
	title    string
	contents []string
}

// 打印章节的统计信息
func statsChapter(ch *chapter) {
	fmt.Printf("Chapter %v has %v lines\n", ch.title, len(ch.contents))
	fmt.Println("Last line:", ch.contents[len(ch.contents)-1:])
}

func main() {
	// 打开文件
	file, err := os.Open("E:/books/gutenberg/pride_and_prejudice.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 使用scanner按行读取
	scanner := bufio.NewScanner(file)
	isContent := false
	var curChapter chapter
	for scanner.Scan() {
		// 当前行
		line := strings.TrimSpace(scanner.Text())

		// 如果遇到正文开始分隔行，则进入正文
		if isStartDelim(line) {
			isContent = true
			continue
		}

		// 如果遇到正文结束分隔行，则退出循环
		if isEndDelim(line) {
			break
		}

		// 如果遇到新章节的标题
		if isChapterTitle(line) {
			// 统计上一章节
			statsChapter(&curChapter)
			// 创建新章节
			chapterNum := extractChapterTitle(line)
			curChapter = chapter{chapterNum, []string{}}
			continue
		}

		// 如果是正文，则当前行加入当前章节数据中
		if isContent && len(line) > 0 {
			curChapter.contents = append(curChapter.contents, line)
		}

	}
	statsChapter(&curChapter)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
